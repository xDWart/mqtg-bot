package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"mqtg-bot/internal/models"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	DEFAULT_POSTGRES_HOST = "localhost"
	DEFAULT_POSTGRES_PORT = "5432"
	DEFAULT_POSTGRES_USER = "postgres"
)

func NewPostgresConnection() *gorm.DB {
	uri, err := url.Parse(os.Getenv("DATABASE_URL"))

	if uri.Host == "" || err != nil {
		host := os.Getenv("POSTGRES_HOST")
		if host == "" {
			host = DEFAULT_POSTGRES_HOST
		}

		port := os.Getenv("POSTGRES_PORT")
		if port == "" {
			port = DEFAULT_POSTGRES_PORT
		}

		uri.Host = host + ":" + port

		user := os.Getenv("POSTGRES_USER")
		if user == "" {
			user = DEFAULT_POSTGRES_USER
		}

		password := os.Getenv("POSTGRES_PASSWORD")
		if password == "" {
			log.Fatalf("POSTGRES_PASSWORD must not be empty or undefined")
		}

		uri.User = url.UserPassword(user, password)

		dbName := os.Getenv("POSTGRES_DB")
		if dbName == "" {
			dbName = user
		}

		uri.Path = "/" + dbName
	}

	var db *gorm.DB

	userPassword, _ := uri.User.Password()
	connectStr := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v", uri.Hostname(), uri.Port(), uri.User.Username(), strings.TrimPrefix(uri.Path, "/"), userPassword)

	sslmode := os.Getenv("DB_SSLMODE")
	if len(sslmode) > 0 {
		connectStr += fmt.Sprintf(" sslmode=%v", sslmode)
	}
	// log.Printf("Connect string: %v", connectStr)

	for i := 0; i < 10; i++ {
		db, err = gorm.Open("postgres", connectStr)
		if err == nil {
			break
		}
		log.Printf("Connect postgres error: %+v", err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		log.Fatalf("Could not init postgres connection: %+v", err)
	}

	if os.Getenv("DB_DEBUG") == "true" {
		db.LogMode(true)
	}

	autoMigrate(db)

	return db
}

func autoMigrate(db *gorm.DB) {
	db.AutoMigrate(
		&models.DbUser{},
		&models.Subscription{},
		&models.SubscriptionData{},
	)
}
