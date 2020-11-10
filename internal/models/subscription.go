package models

import (
	"github.com/jinzhu/gorm"
	"sync"
	"time"
)

type Subscription struct {
	gorm.Model
	UserMutex                 *sync.Mutex `gorm:"-"`
	DbUserID                  uint
	ChatID                    int64                `gorm:"type:bigint"`
	Topic                     string               `gorm:"type:varchar(255)"`
	Qos                       byte                 `gorm:"default:0"`
	DataType                  SubscriptionDataType `gorm:"default:0"`
	SubscriptionType          SubscriptionType     `gorm:"default:0"`
	CollectedData             []SubscriptionData
	BeforeValueText           string `gorm:"type:text"`
	AfterValueText            string `gorm:"type:text"`
	LastValueFormattedMessage string `gorm:"type:text"`
	LastValuePayload          []byte `gorm:"type:bytea"`
	JsonPathToData            string `gorm:"type:varchar(255)"`
}

type SubscriptionData struct {
	gorm.Model
	SubscriptionID uint
	DateTime       time.Time
	DataType       SubscriptionDataType `gorm:"default:0"`
	Data           []byte               `gorm:"type:bytea"`
}

type SubscriptionDataType byte

const (
	TEXT_DATA_TYPE SubscriptionDataType = iota
	IMAGE_DATA_TYPE
	COUNT_DATA_TYPES
)

var dataTypeStrings = [...]string{
	TEXT_DATA_TYPE:  "text",
	IMAGE_DATA_TYPE: "image",
}

var _ = dataTypeStrings[COUNT_DATA_TYPES-1]

func (dt SubscriptionDataType) String() string {
	return dataTypeStrings[byte(dt)]
}

func (dt SubscriptionDataType) GetNext() SubscriptionDataType {
	return (dt + 1) % COUNT_DATA_TYPES
}

type SubscriptionType byte

const (
	PRINT_MESSAGE_WITHOUT_STORING_SUBSCRIPTION_TYPE SubscriptionType = iota
	PRINT_AND_STORE_MESSAGE_SUBSCRIPTION_TYPE
	SILENT_STORE_MESSAGE_SUBSCRIPTION_TYPE
	COUNT_SUBSCRIPTION_TYPES
)

var subscriptionTypeStrings = [...]string{
	PRINT_MESSAGE_WITHOUT_STORING_SUBSCRIPTION_TYPE: "print message without storing",
	PRINT_AND_STORE_MESSAGE_SUBSCRIPTION_TYPE:       "print and store message",
	SILENT_STORE_MESSAGE_SUBSCRIPTION_TYPE:          "silent store message",
}

var _ = subscriptionTypeStrings[COUNT_SUBSCRIPTION_TYPES-1]

func (st SubscriptionType) String() string {
	return subscriptionTypeStrings[byte(st)]
}

func (st SubscriptionType) GetNext() SubscriptionType {
	return (st + 1) % COUNT_SUBSCRIPTION_TYPES
}
