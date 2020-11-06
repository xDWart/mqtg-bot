package button_interface

import (
	"encoding/json"
	"mqtg-bot/internal/users/menu/button_types"
)

var _ ButtonI = &PrintLastValueButton{}

type PrintLastValueButton struct {
	Parent         *FolderButton `json:"-"`
	Type           button_types.ButtonType
	Name           string
	SubscriptionID int `json:"SubscriptionID"`
}

func (b *PrintLastValueButton) GetType() button_types.ButtonType {
	return button_types.PRINT_LAST_SUB_VALUE
}
func (b *PrintLastValueButton) GetName() string {
	return b.Name
}
func (b *PrintLastValueButton) GetFullName() string {
	return b.Name
}
func (b *PrintLastValueButton) SetMainName(name string) {
	b.Name = name
}

func (b *PrintLastValueButton) GetCurrentCommand() *CommandType {
	return nil
}
func (b *PrintLastValueButton) GetCommands() []*CommandType {
	return nil
}
func (b *PrintLastValueButton) AddNewCommand(*CommandType) {}
func (b *PrintLastValueButton) DeleteCommand(int)          {}

func (b *PrintLastValueButton) SwitchState()                    {}
func (b *PrintLastValueButton) SetNameForCommand(int, string)   {}
func (b *PrintLastValueButton) SetTopicForCommand(int, string)  {}
func (b *PrintLastValueButton) SetValueForCommand(int, string)  {}
func (b *PrintLastValueButton) SetQosForCommand(int, byte)      {}
func (b *PrintLastValueButton) SetRetainedForCommand(int, bool) {}

func (b *PrintLastValueButton) GetSubscriptions() []int {
	return []int{b.SubscriptionID}
}
func (b *PrintLastValueButton) SetSubscription(i int, s int) {
	b.SubscriptionID = s
}

func (b *PrintLastValueButton) SetParent(parent *FolderButton) {
	b.Parent = parent
}

func (b *PrintLastValueButton) GetParent() *FolderButton {
	return b.Parent
}

func (b *PrintLastValueButton) GetButtons() *[]ButtonI {
	return nil
}

func (b *PrintLastValueButton) AddButton(ButtonI) {
}

func (b *PrintLastValueButton) DelButton(int32) {
}

func (b *PrintLastValueButton) MarshalJSON() ([]byte, error) {
	b.Type = b.GetType()
	return json.Marshal(*b)
}

func (b *PrintLastValueButton) UnmarshalJSON([]byte) error {
	return nil
}
