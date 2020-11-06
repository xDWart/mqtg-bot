package internal

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"mqtg-bot/internal/database"
	"mqtg-bot/internal/users"
	"mqtg-bot/internal/users/mqtt"
	"os"
	"runtime"
	"strconv"
	"sync"
)

type TelegramBot struct {
	*tgbotapi.BotAPI
	db             *gorm.DB
	updates        tgbotapi.UpdatesChannel
	subscriptionCh chan mqtt.SubscriptionMessage
	usersManager   *users.Manager

	wg              *sync.WaitGroup
	shutdownChannel chan interface{}
	metrics         Metrics
}

func InitTelegramBot() *TelegramBot {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatalf("TELEGRAM_BOT_TOKEN does not set")
	}

	botApi, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf(err.Error())
	}

	bot := &TelegramBot{
		BotAPI:          botApi,
		db:              database.NewPostgresConnection(),
		subscriptionCh:  make(chan mqtt.SubscriptionMessage),
		metrics:         InitPrometheusMetrics(),
		wg:              &sync.WaitGroup{},
		shutdownChannel: make(chan interface{}),
	}

	if os.Getenv("BOT_DEBUG") == "true" {
		bot.Debug = true
	}

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	bot.updates = bot.GetUpdatesChan(updateConfig)
	bot.usersManager = users.InitManager(bot.db, bot.subscriptionCh)
	bot.usersManager.LoadAllConnectedUsers()

	prometheus.MustRegister(bot.metrics.GetPrometheusMetrics()...)
	prometheus.MustRegister(bot.usersManager.GetPrometheusMetrics()...)
	prometheus.MustRegister(mqtt.GetPrometheusMetrics()...)
	go bot.StartPprofAndMetricsListener()

	log.Printf("Successfully init Telegram Bot")

	var numListenGoroutines int
	if os.Getenv("NUM_LISTEN_GOROUTINES") != "" {
		numListenGoroutines, err = strconv.Atoi(os.Getenv("NUM_LISTEN_GOROUTINES"))
		if err != nil {
			log.Printf("Bad NUM_LISTEN_GOROUTINES env: %v", err)
		}
	}

	if numListenGoroutines == 0 {
		numListenGoroutines = runtime.NumCPU()
		log.Printf("NUM_LISTEN_GOROUTINES is not set, by default will use NumCPU(%v) goroutines", numListenGoroutines)
	}

	log.Printf("Running %v listeners", numListenGoroutines)

	for i := 0; i < numListenGoroutines; i++ {
		bot.wg.Add(1)
		go bot.StartBotListener()
	}

	return bot
}

func (bot *TelegramBot) Shutdown() {
	log.Printf("Telegram Bot received shutdown signal, will close all listeners")

	bot.StopReceivingUpdates()
	close(bot.shutdownChannel)
	bot.db.Close()
}

func (bot *TelegramBot) Wait() {
	bot.wg.Wait()
}
