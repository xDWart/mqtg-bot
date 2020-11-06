package button_interface

import (
	"encoding/json"
	"mqtg-bot/internal/users/menu/button_types"
)

var _ ButtonI = &ToggleButton{}

type ToggleButton struct {
	Parent   *FolderButton `json:"-"`
	Type     button_types.ButtonType
	State    int
	Commands []*CommandType
}

func (b *ToggleButton) GetType() button_types.ButtonType {
	return button_types.TOGGLE
}

func (b *ToggleButton) GetName() string {
	extendCommandsSliceIfNeeded(b.State, &b.Commands)
	return b.Commands[b.State].Name
}
func (b *ToggleButton) GetFullName() string {
	extendCommandsSliceIfNeeded(1, &b.Commands)
	return b.Commands[0].Name + " / " + b.Commands[1].Name
}
func (b *ToggleButton) SetMainName(name string) {
	extendCommandsSliceIfNeeded(b.State, &b.Commands)
	b.Commands[b.State].Name = name
}

func (b *ToggleButton) GetCurrentCommand() *CommandType {
	extendCommandsSliceIfNeeded(b.State, &b.Commands)
	return b.Commands[b.State]
}
func (b *ToggleButton) GetCommands() []*CommandType {
	return b.Commands
}
func (b *ToggleButton) AddNewCommand(*CommandType) {}
func (b *ToggleButton) DeleteCommand(int)          {}

func (b *ToggleButton) SwitchState() {
	b.State = (b.State + 1) % 2
}

func (b *ToggleButton) SetNameForCommand(s int, name string) {
	extendCommandsSliceIfNeeded(s, &b.Commands)
	b.Commands[s].Name = name
}
func (b *ToggleButton) SetTopicForCommand(s int, topic string) {
	extendCommandsSliceIfNeeded(s, &b.Commands)
	b.Commands[s].Topic = topic
}
func (b *ToggleButton) SetValueForCommand(s int, value string) {
	extendCommandsSliceIfNeeded(s, &b.Commands)
	b.Commands[s].Value = value
}
func (b *ToggleButton) SetQosForCommand(s int, qos byte) {
	extendCommandsSliceIfNeeded(s, &b.Commands)
	b.Commands[s].Qos = qos
}
func (b *ToggleButton) SetRetainedForCommand(s int, retained bool) {
	extendCommandsSliceIfNeeded(s, &b.Commands)
	b.Commands[s].Retained = retained
}

func (b *ToggleButton) GetSubscriptions() []int {
	return nil
}
func (b *ToggleButton) SetSubscription(int, int) {}

func (b *ToggleButton) SetParent(parent *FolderButton) {
	b.Parent = parent
}

func (b *ToggleButton) GetParent() *FolderButton {
	return b.Parent
}

func (b *ToggleButton) GetButtons() *[]ButtonI {
	return nil
}

func (b *ToggleButton) AddButton(ButtonI) {
}

func (b *ToggleButton) DelButton(int32) {
}

func (b *ToggleButton) MarshalJSON() ([]byte, error) {
	b.Type = b.GetType()
	return json.Marshal(*b)
}

func (b *ToggleButton) UnmarshalJSON([]byte) error {
	return nil
}
