package models

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type DbUser struct {
	gorm.Model
	ChatID        int64  `gorm:"type:bigint"`
	UserName      string `gorm:"type:varchar(255)"`
	MqttUrl       string `gorm:"type:varchar(255)"` // (tcp|ssl|ws|wss)://user:password@host:port/path
	Connected     bool
	Subscriptions []*Subscription
	DbMenu        postgres.Jsonb
}
