package state

import (
	"mqtg-bot/internal/models"
	"mqtg-bot/internal/users/menu/button_interface"
	"mqtg-bot/internal/users/menu/button_types"
)

type StateStruct struct {
	State                StateType
	LastMessageID        int
	CurrPath             []int32
	EditableIndex        int32
	PublishingTopic      string
	Qos                  byte
	Retained             bool
	ButtonType           button_types.ButtonType
	SubscriptionType     models.SubscriptionType
	SubscriptionDataType models.SubscriptionDataType
	EditableButton       button_interface.ButtonI
}

type StateType byte

const (
	NIL_STATE StateType = iota
	CONFIGURE_CONNECTION_STATE
	PUBLISH_TOPIC_STATE
	PUBLISH_VALUE_STATE
	ADD_SUBSCRIPTION_STATE
	EDIT_BEFORE_VALUE_MESSAGE_TEXT_STATE
	EDIT_AFTER_VALUE_MESSAGE_TEXT_STATE
	EDIT_SUBSCRIPTION_TOPIC_STATE
	ADD_BUTTON_STATE
	EDIT_BUTTON_NAME_STATE
	ADD_NEW_COMMAND
	EDIT_COMMAND_NAME_STATE
	EDIT_COMMAND_TOPIC_STATE
	EDIT_COMMAND_VALUE_STATE
)

func (state *StateStruct) Reset() {
	state.State = NIL_STATE
	state.LastMessageID = 0
	state.CurrPath = []int32{}
	state.EditableIndex = 0
	state.PublishingTopic = ""
	state.Qos = 0
	state.Retained = false
	state.ButtonType = 0
	state.SubscriptionType = 0
	state.SubscriptionDataType = 0
	state.EditableButton = nil
}
