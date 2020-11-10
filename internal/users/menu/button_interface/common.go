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

		name, _ := (*dataMap)["Name"].(string)
		multiValueButton.Name = name

		commands, _ := (*dataMap)["Commands"].([]interface{})
		multiValueButton.Commands = make([]*CommandType, 0, len(commands))
		for _, command := range commands {
			commandI, ok := command.(map[string]interface{})
			if ok {
				var commandT CommandType
				commandT.Name, _ = commandI["Name"].(string)
				commandT.Topic, _ = commandI["Topic"].(string)
				commandT.Value, _ = commandI["Value"].(string)
				commandT.Qos, _ = commandI["Qos"].(byte)
				commandT.Retained, _ = commandI["Retained"].(bool)
				multiValueButton.Commands = append(multiValueButton.Commands, &commandT)
			}
		}

		outButton = &multiValueButton

	case button_types.SINGLE_VALUE:
		var singleValueButton SingleValueButton

		command, ok := (*dataMap)["Command"].(map[string]interface{})
		if ok {
			singleValueButton.Command.Name, _ = command["Name"].(string)
			singleValueButton.Command.Topic, _ = command["Topic"].(string)
			singleValueButton.Command.Value, _ = command["Value"].(string)
			singleValueButton.Command.Qos, _ = command["Qos"].(byte)
			singleValueButton.Command.Retained, _ = command["Retained"].(bool)
		}
		outButton = &singleValueButton

	case button_types.TOGGLE:
		var toggleButton ToggleButton
		state, _ := (*dataMap)["State"].(float64)
		toggleButton.State = int(state)

		commands, _ := (*dataMap)["Commands"].([]interface{})

		toggleButton.Commands = make([]*CommandType, 0, len(commands))
		for _, command := range commands {
			commandI, ok := command.(map[string]interface{})
			if ok {
				var commandT CommandType
				commandT.Name, _ = commandI["Name"].(string)
				commandT.Topic, _ = commandI["Topic"].(string)
				commandT.Value, _ = commandI["Value"].(string)
				commandT.Qos, _ = commandI["Qos"].(byte)
				commandT.Retained, _ = commandI["Retained"].(bool)
				toggleButton.Commands = append(toggleButton.Commands, &commandT)
			}
		}

		outButton = &toggleButton

	case button_types.FOLDER:
		var folderButton FolderButton
		folderButton.Name, _ = (*dataMap)["Name"].(string)
		buttons, _ := (*dataMap)["Buttons"].([]interface{})

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
		systemButton.Name, _ = (*dataMap)["Name"].(string)
		outButton = &systemButton

	case button_types.PRINT_LAST_SUB_VALUE:
		var printLastValueButton PrintLastValueButton
		printLastValueButton.Name, _ = (*dataMap)["Name"].(string)
		subscriptionID, _ := (*dataMap)["SubscriptionID"].(float64)
		printLastValueButton.SubscriptionID = int(subscriptionID)
		outButton = &printLastValueButton

	case button_types.DRAW_CHART:
		var showCharButton DrawChartButton
		showCharButton.Name, _ = (*dataMap)["Name"].(string)

		subscriptions, _ := (*dataMap)["Subscriptions"].([]interface{})
		showCharButton.Subscriptions = make([]int, 0, len(subscriptions))
		for _, subscription := range subscriptions {
			subscriptionID, _ := subscription.(float64)
			showCharButton.Subscriptions = append(showCharButton.Subscriptions, int(subscriptionID))
		}

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
