package button_interface

import (
	"fmt"
	"log"
	"mqtg-bot/internal/users/menu/button_types"
)

type CommandType struct {
	Name     string
	Topic    string
	Value    string
	Qos      byte
	Retained bool
}

func GetNewButtonWithName(buttonType button_types.ButtonType, buttonName string) (ButtonI, error) {
	switch buttonType {
	case button_types.FOLDER:
		return &FolderButton{
			Name: buttonName,
		}, nil
	case button_types.SINGLE_VALUE:
		return &SingleValueButton{
			Command: CommandType{
				Name: buttonName,
			},
		}, nil
	case button_types.TOGGLE:
		return &ToggleButton{
			Commands: []*CommandType{
				{
					Name: buttonName,
				},
				{
					Name: buttonName,
				},
			},
		}, nil
	case button_types.MULTI_VALUE:
		return &MultiValueButton{
			Name:     buttonName,
			Commands: []*CommandType{},
		}, nil
	case button_types.PRINT_LAST_SUB_VALUE:
		return &PrintLastValueButton{
			Name:           buttonName,
			SubscriptionID: -1,
		}, nil
	case button_types.DRAW_CHART:
		return &DrawChartButton{
			Name:          buttonName,
			Subscriptions: []int{},
		}, nil
	}
	return nil, fmt.Errorf("unknown button type: %v", buttonType)
}

func parseDataMap(dataMap *map[string]interface{}) ButtonI {
	var outButton ButtonI

	fType, _ := (*dataMap)["Type"].(float64)

	switch button_types.ButtonType(fType) {
	case button_types.MULTI_VALUE:
		var multiValueButton MultiValueButton
		// заполнение нужных полей
		outButton = &multiValueButton

	case button_types.SINGLE_VALUE:
		var singleValueButton SingleValueButton
		// заполнение нужных полей
		outButton = &singleValueButton

	case button_types.TOGGLE:
		var toggleButton ToggleButton
		// заполнение нужных полей
		outButton = &toggleButton

	case button_types.FOLDER:
		var folderButton FolderButton
		folderButton.Name, _ = (*dataMap)["Name"].(string)
		buttons, _ := (*dataMap)["Buttons"].([]interface{})

		// рекурсивный вызов парсинга для дерева
		folderButton.Buttons = make([]ButtonI, 0, len(buttons))
		for _, button := range buttons {
			buttonMap, ok := button.(map[string]interface{})
			if ok {
				ButtonI := parseDataMap(&buttonMap)
				folderButton.Buttons = append(folderButton.Buttons, ButtonI)
			}
		}

		outButton = &folderButton

	case button_types.SYSTEM:
		var systemButton SystemButton
		// заполнение нужных полей
		outButton = &systemButton

	case button_types.PRINT_LAST_SUB_VALUE:
		var printLastValueButton PrintLastValueButton
		// заполнение нужных полей
		outButton = &printLastValueButton

	case button_types.DRAW_CHART:
		var showCharButton DrawChartButton
		// заполнение нужных полей
		outButton = &showCharButton

	default:
		log.Printf("Unknown button Type: %#v", *dataMap)
	}

	return outButton
}

func extendCommandsSliceIfNeeded(s int, commands *[]*CommandType) {
	if s >= len(*commands) { // need add new command
		comCount := len(*commands)
		for i := 0; i <= s-comCount; i++ {
			*commands = append(*commands, &CommandType{})
		}
	}
}
