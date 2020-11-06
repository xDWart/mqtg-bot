package button_interface

import (
	"encoding/json"
	"mqtg-bot/internal/users/menu/button_types"
)

var _ ButtonI = &DrawChartButton{}

type DrawChartButton struct {
	Parent        *FolderButton `json:"-"`
	Type          button_types.ButtonType
	Name          string
	Subscriptions []int
}

func (b *DrawChartButton) GetType() button_types.ButtonType {
	return button_types.DRAW_CHART
}
func (b *DrawChartButton) GetName() string {
	return b.Name
}
func (b *DrawChartButton) GetFullName() string {
	return b.Name
}
func (b *DrawChartButton) SetMainName(name string) {
	b.Name = name
}

func (b *DrawChartButton) GetCurrentCommand() *CommandType {
	return nil
}
func (b *DrawChartButton) GetCommands() []*CommandType {
	return nil
}
func (b *DrawChartButton) AddNewCommand(*CommandType) {}
func (b *DrawChartButton) DeleteCommand(int)          {}

func (b *DrawChartButton) SwitchState()                    {}
func (b *DrawChartButton) SetNameForCommand(int, string)   {}
func (b *DrawChartButton) SetTopicForCommand(int, string)  {}
func (b *DrawChartButton) SetValueForCommand(int, string)  {}
func (b *DrawChartButton) SetQosForCommand(int, byte)      {}
func (b *DrawChartButton) SetRetainedForCommand(int, bool) {}

func (b *DrawChartButton) GetSubscriptions() []int {
	return b.Subscriptions
}
func (b *DrawChartButton) SetSubscription(k int, s int) {
	if k >= len(b.Subscriptions) {
		comCount := len(b.Subscriptions)
		for i := 0; i <= k-comCount; i++ {
			b.Subscriptions = append(b.Subscriptions, 0)
		}
	}
	if s < 0 { // delete
		b.Subscriptions = append(b.Subscriptions[:k], b.Subscriptions[k+1:]...)
	} else {
		b.Subscriptions[k] = s
	}
}

func (b *DrawChartButton) SetParent(parent *FolderButton) {
	b.Parent = parent
}

func (b *DrawChartButton) GetParent() *FolderButton {
	return b.Parent
}

func (b *DrawChartButton) GetButtons() *[]ButtonI {
	return nil
}

func (b *DrawChartButton) AddButton(ButtonI) {
}

func (b *DrawChartButton) DelButton(int32) {
}

func (b *DrawChartButton) MarshalJSON() ([]byte, error) {
	b.Type = b.GetType()
	return json.Marshal(*b)
}

func (b *DrawChartButton) UnmarshalJSON([]byte) error {
	return nil
}
