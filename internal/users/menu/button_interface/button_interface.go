package button_interface

import "mqtg-bot/internal/users/menu/button_types"

type ButtonI interface {
	GetType() button_types.ButtonType

	GetName() string
	GetFullName() string // for TOGGLE button only
	SetMainName(string)

	GetCurrentCommand() *CommandType
	GetCommands() []*CommandType

	AddNewCommand(*CommandType)
	DeleteCommand(int)

	SwitchState() // for TOGGLE button only
	SetNameForCommand(int, string)
	SetTopicForCommand(int, string)
	SetValueForCommand(int, string)
	SetQosForCommand(int, byte)
	SetRetainedForCommand(int, bool)

	GetSubscriptions() []int  // for PRINT_LAST_SUB_VALUE and DRAW_CHART buttons only
	SetSubscription(int, int) // for PRINT_LAST_SUB_VALUE and DRAW_CHART buttons only

	SetParent(*FolderButton)
	GetParent() *FolderButton

	GetButtons() *[]ButtonI
	AddButton(ButtonI)
	DelButton(int32)

	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
}
