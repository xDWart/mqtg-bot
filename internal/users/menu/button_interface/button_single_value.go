package button_interface

import (
	"encoding/json"
	"mqtg-bot/internal/users/menu/button_types"
)

var _ ButtonI = &SingleValueButton{}

type SingleValueButton struct {
	Parent  *FolderButton `json:"-"`
	Type    button_types.ButtonType
	Command CommandType
}

func (b *SingleValueButton) GetType() button_types.ButtonType {
	return button_types.SINGLE_VALUE
}

func (b *SingleValueButton) GetName() string {
	return b.Command.Name
}
func (b *SingleValueButton) GetFullName() string {
	return b.Command.Name
}
func (b *SingleValueButton) SetMainName(name string) {
	b.Command.Name = name
}

func (b *SingleValueButton) GetCurrentCommand() *CommandType {
	return &b.Command
}
func (b *SingleValueButton) GetCommands() []*CommandType {
	return []*CommandType{&b.Command}
}
func (b *SingleValueButton) AddNewCommand(*CommandType) {}
func (b *SingleValueButton) DeleteCommand(int)          {}

func (b *SingleValueButton) SwitchState() {}

func (b *SingleValueButton) SetNameForCommand(s int, name string) {
	b.Command.Name = name
}
func (b *SingleValueButton) SetTopicForCommand(s int, topic string) {
	b.Command.Topic = topic
}
func (b *SingleValueButton) SetValueForCommand(s int, value string) {
	b.Command.Value = value
}
func (b *SingleValueButton) SetQosForCommand(s int, qos byte) {
	b.Command.Qos = qos
}
func (b *SingleValueButton) SetRetainedForCommand(s int, retained bool) {
	b.Command.Retained = retained
}

func (b *SingleValueButton) GetSubscriptions() []int {
	return nil
}
func (b *SingleValueButton) SetSubscription(int, int) {}

func (b *SingleValueButton) SetParent(parent *FolderButton) {
	b.Parent = parent
}

func (b *SingleValueButton) GetParent() *FolderButton {
	return b.Parent
}

func (b *SingleValueButton) GetButtons() *[]ButtonI {
	return nil
}

func (b *SingleValueButton) AddButton(ButtonI) {
}

func (b *SingleValueButton) DelButton(int32) {
}

func (b *SingleValueButton) MarshalJSON() ([]byte, error) {
	b.Type = b.GetType()
	return json.Marshal(*b)
}

func (b *SingleValueButton) UnmarshalJSON([]byte) error {
	return nil
}
