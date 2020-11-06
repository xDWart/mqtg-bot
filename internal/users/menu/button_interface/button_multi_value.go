package button_interface

import (
	"encoding/json"
	"mqtg-bot/internal/users/menu/button_types"
)

var _ ButtonI = &MultiValueButton{}

type MultiValueButton struct {
	Parent   *FolderButton `json:"-"`
	Type     button_types.ButtonType
	Name     string
	Commands []*CommandType
}

func (b *MultiValueButton) GetType() button_types.ButtonType {
	return button_types.MULTI_VALUE
}

func (b *MultiValueButton) GetName() string {
	return b.Name
}
func (b *MultiValueButton) GetFullName() string {
	return b.Name
}
func (b *MultiValueButton) SetMainName(name string) {
	b.Name = name
}

func (b *MultiValueButton) GetCurrentCommand() *CommandType {
	return nil
}
func (b *MultiValueButton) GetCommands() []*CommandType {
	return b.Commands
}

func (b *MultiValueButton) AddNewCommand(newCommand *CommandType) {
	b.Commands = append(b.Commands, newCommand)
}
func (b *MultiValueButton) DeleteCommand(s int) {
	b.Commands = append(b.Commands[:s], b.Commands[s+1:]...)
}

func (b *MultiValueButton) SwitchState() {}

func (b *MultiValueButton) SetNameForCommand(s int, name string) {
	extendCommandsSliceIfNeeded(s, &b.Commands)
	b.Commands[s].Name = name
}
func (b *MultiValueButton) SetTopicForCommand(s int, topic string) {
	extendCommandsSliceIfNeeded(s, &b.Commands)
	b.Commands[s].Topic = topic
}
func (b *MultiValueButton) SetValueForCommand(s int, value string) {
	extendCommandsSliceIfNeeded(s, &b.Commands)
	b.Commands[s].Value = value
}
func (b *MultiValueButton) SetQosForCommand(s int, qos byte) {
	extendCommandsSliceIfNeeded(s, &b.Commands)
	b.Commands[s].Qos = qos
}
func (b *MultiValueButton) SetRetainedForCommand(s int, retained bool) {
	extendCommandsSliceIfNeeded(s, &b.Commands)
	b.Commands[s].Retained = retained
}

func (b *MultiValueButton) GetSubscriptions() []int {
	return nil
}
func (b *MultiValueButton) SetSubscription(int, int) {}

func (b *MultiValueButton) SetParent(parent *FolderButton) {
	b.Parent = parent
}

func (b *MultiValueButton) GetParent() *FolderButton {
	return b.Parent
}

func (b *MultiValueButton) GetButtons() *[]ButtonI {
	return nil
}

func (b *MultiValueButton) AddButton(ButtonI) {
}

func (b *MultiValueButton) DelButton(int32) {
}

func (b *MultiValueButton) MarshalJSON() ([]byte, error) {
	b.Type = b.GetType()
	return json.Marshal(*b)
}

func (b *MultiValueButton) UnmarshalJSON([]byte) error {
	return nil
}
