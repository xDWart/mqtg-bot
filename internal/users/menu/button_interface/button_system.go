package button_interface

import (
	"encoding/json"
	"mqtg-bot/internal/users/menu/button_types"
)

var _ ButtonI = &SystemButton{}

type SystemButton struct {
	Parent *FolderButton `json:"-"`
	Type   button_types.ButtonType
	Name   string
}

func (b *SystemButton) GetType() button_types.ButtonType {
	return button_types.SYSTEM
}
func (b *SystemButton) GetName() string {
	return b.Name
}
func (b *SystemButton) GetFullName() string {
	return b.Name
}
func (b *SystemButton) SetMainName(name string) {
	b.Name = name
}

func (b *SystemButton) GetCurrentCommand() *CommandType {
	return nil
}
func (b *SystemButton) GetCommands() []*CommandType {
	return nil
}
func (b *SystemButton) AddNewCommand(*CommandType) {}
func (b *SystemButton) DeleteCommand(int)          {}

func (b *SystemButton) SwitchState()                    {}
func (b *SystemButton) SetNameForCommand(int, string)   {}
func (b *SystemButton) SetTopicForCommand(int, string)  {}
func (b *SystemButton) SetValueForCommand(int, string)  {}
func (b *SystemButton) SetQosForCommand(int, byte)      {}
func (b *SystemButton) SetRetainedForCommand(int, bool) {}

func (b *SystemButton) GetSubscriptions() []int {
	return nil
}
func (b *SystemButton) SetSubscription(int, int) {}

func (b *SystemButton) SetParent(parent *FolderButton) {
	b.Parent = parent
}

func (b *SystemButton) GetParent() *FolderButton {
	return b.Parent
}

func (b *SystemButton) GetButtons() *[]ButtonI {
	return nil
}

func (b *SystemButton) AddButton(ButtonI) {
}

func (b *SystemButton) DelButton(int32) {
}

func (b *SystemButton) MarshalJSON() ([]byte, error) {
	b.Type = b.GetType()
	return json.Marshal(*b)
}

func (b *SystemButton) UnmarshalJSON([]byte) error {
	return nil
}
