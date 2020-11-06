package button_interface

import (
	"encoding/json"
	"mqtg-bot/internal/users/menu/button_types"
)

var _ ButtonI = &FolderButton{}

type FolderButton struct {
	Parent  *FolderButton `json:"-"`
	Type    button_types.ButtonType
	Name    string
	Buttons []ButtonI
}

func (b *FolderButton) GetType() button_types.ButtonType {
	return button_types.FOLDER
}

func (b *FolderButton) GetName() string {
	return b.Name
}
func (b *FolderButton) GetFullName() string {
	return b.Name
}
func (b *FolderButton) SetMainName(name string) {
	b.Name = name
}

func (b *FolderButton) GetCurrentCommand() *CommandType {
	return nil
}
func (b *FolderButton) GetCommands() []*CommandType {
	return nil
}

func (b *FolderButton) AddNewCommand(*CommandType) {}
func (b *FolderButton) DeleteCommand(int)          {}

func (b *FolderButton) SwitchState()                    {}
func (b *FolderButton) SetNameForCommand(int, string)   {}
func (b *FolderButton) SetTopicForCommand(int, string)  {}
func (b *FolderButton) SetValueForCommand(int, string)  {}
func (b *FolderButton) SetQosForCommand(int, byte)      {}
func (b *FolderButton) SetRetainedForCommand(int, bool) {}

func (b *FolderButton) GetSubscriptions() []int {
	return nil
}
func (b *FolderButton) SetSubscription(int, int) {}

func (b *FolderButton) SetParent(parent *FolderButton) {
	b.Parent = parent
	for i := range b.Buttons {
		b.Buttons[i].SetParent(b)
	}
}

func (b *FolderButton) GetParent() *FolderButton {
	return b.Parent
}

func (b *FolderButton) GetButtons() *[]ButtonI {
	return &b.Buttons
}

func (b *FolderButton) AddButton(button ButtonI) {
	b.Buttons = append(b.Buttons, button)
	button.SetParent(b)
}

func (b *FolderButton) DelButton(i int32) {
	b.Buttons = append(b.Buttons[:i], b.Buttons[i+1:]...)
}

func (b *FolderButton) MarshalJSON() ([]byte, error) {
	b.Type = b.GetType()
	return json.Marshal(*b)
}

func (b *FolderButton) UnmarshalJSON(data []byte) error {
	var dataMap map[string]interface{}

	err := json.Unmarshal(data, &dataMap)
	if err != nil {
		return err
	}

	folderButton := parseDataMap(&dataMap).(*FolderButton)
	b.Name = folderButton.Name
	b.Buttons = folderButton.Buttons

	return nil
}
