package users

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"mqtg-bot/internal/models"
	"mqtg-bot/internal/users/menu"
	"mqtg-bot/internal/users/mqtt"
	"sync"
)

type Manager struct {
	sync.RWMutex
	db             *gorm.DB
	subscriptionCh chan mqtt.SubscriptionMessage
	users          map[int64]*User
	metrics        Metrics
}

func InitManager(db *gorm.DB, subscriptionCh chan mqtt.SubscriptionMessage) *Manager {
	return &Manager{
		db:             db,
		subscriptionCh: subscriptionCh,
		users:          make(map[int64]*User),
		metrics:        InitPrometheusMetrics(),
	}
}

func (um *Manager) LoadAllConnectedUsers() {
	// select all users connected to MQTT
	var dbUsers []*models.DbUser
	um.db.Where("connected = ?", true).Find(&dbUsers)

	// convert dbUser into UserType
	for _, dbUser := range dbUsers {
		log.Printf("Loaded connected user %v", dbUser.UserName)
		um.LoadDatabaseUserIntoBotUsers(dbUser)
	}

	um.UpdateTotalUsers()
}

func (um *Manager) UpdateTotalUsers() {
	var count float64
	um.db.Model(&models.DbUser{}).Count(&count)
	um.metrics.numOfTotalUsers.Set(count)
}

func (um *Manager) GetUserByChatIdFromUpdate(update *tgbotapi.Update) *User {
	var message = update.Message
	if message == nil {
		if update.CallbackQuery != nil {
			message = update.CallbackQuery.Message
		} else {
			return nil
		}
	}

	um.RLock()
	botUser, ok := um.users[message.Chat.ID]
	um.RUnlock()

	if !ok { // user not found, need to create
		var dbUser models.DbUser

		// first try select from db
		um.db.Where("chat_id = ?", message.Chat.ID).First(&dbUser)
		if dbUser.ID == 0 {
			// create a new one
			dbUser.ChatID = message.Chat.ID
			dbUser.UserName = message.From.UserName

			um.db.Create(&dbUser)

			um.UpdateTotalUsers()
		}

		botUser = um.LoadDatabaseUserIntoBotUsers(&dbUser)
	}

	return botUser
}

func (um *Manager) LoadDatabaseUserIntoBotUsers(dbUser *models.DbUser) *User {
	um.db.Model(dbUser).Order("subscriptions.id").Related(&dbUser.Subscriptions)

	botUser := &User{
		DbUser:         dbUser,
		db:             um.db,
		subscriptionCh: um.subscriptionCh,
		menu:           &menu.MainMenu{},
	}

	if len(dbUser.DbMenu.RawMessage) > 0 {
		botUser.menu.LoadMenuFromJsonb(dbUser.DbMenu)
	}

	botUser.menu.AppendCommonMenuAndSetParentLinks()

	if botUser.Connected { // need to connect
		err := botUser.connectMqttAndSubscribe()
		if err != nil {
			log.Printf("User %v connect MQTT error: %v", botUser.UserName, err)
			botUser.setConnected(false)
		}
	}

	um.Lock()
	um.users[botUser.ChatID] = botUser
	um.metrics.numOfActiveUsers.Set(float64(len(um.users)))
	um.Unlock()

	return botUser
}

func (um *Manager) GetPrometheusMetrics() []prometheus.Collector {
	return um.metrics.GetPrometheusMetrics()
}
