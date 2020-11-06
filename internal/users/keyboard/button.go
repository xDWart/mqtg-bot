package keyboard

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"mqtg-bot/internal/models"
	"mqtg-bot/internal/users/keyboard/callback_data"
	"mqtg-bot/internal/users/keyboard/keyboard_names"
	"mqtg-bot/internal/users/menu/button_interface"
	"mqtg-bot/internal/users/menu/button_names"
	"mqtg-bot/internal/users/menu/button_types"
)

func GetButtonsKeyboard(buttonI button_interface.ButtonI, callbackDataPath []int32, userSubscriptions []*models.Subscription) (string, *tgbotapi.InlineKeyboardMarkup) {
	var text string
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup()

	switch buttonI.GetType() {
	case button_types.FOLDER:
		if len(callbackDataPath) > 0 {
			inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard,
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(
						keyboard_names.RENAME_FOLDER,
						callback_data.QueryDataType{
							Keyboard: callback_data.KeyboardType_BUTTONS,
							Path:     callbackDataPath,
							Action:   callback_data.ActionType_EDIT,
						}.GetBase64ProtoString()),
				),
			)
		}

		for index, button := range *buttonI.GetButtons() {
			inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf("%v: %v", button.GetType(), button.GetFullName()),
					callback_data.QueryDataType{
						Keyboard: callback_data.KeyboardType_BUTTONS,
						Path:     append(callbackDataPath, int32(index)),
					}.GetBase64ProtoString()),
			))
		}

		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.ADD_BUTTON,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_BUTTONS,
					Path:     callbackDataPath,
					Action:   callback_data.ActionType_ADD_BUTTON,
				}.GetBase64ProtoString()),
		))

		// cannot edit main menu
		if len(callbackDataPath) > 0 {
			// can delete only if empty
			if len(*buttonI.GetButtons()) == 0 {
				inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(
						keyboard_names.DELETE_FOLDER,
						callback_data.QueryDataType{
							Keyboard: callback_data.KeyboardType_BUTTONS,
							Path:     callbackDataPath,
							Action:   callback_data.ActionType_DELETE,
						}.GetBase64ProtoString()),
				))
			}

			text = fmt.Sprintf("%v <code>%v</code>:", buttonI.GetType(), buttonI.GetName())
		} else { // top level
			text = "Your buttons menu:"
			if len(*buttonI.GetButtons()) == 0 {
				text += " empty"
			}
		}
	case button_types.SINGLE_VALUE:
		text = fmt.Sprintf("Edit %v <code>%v</code>:", buttonI.GetType(), buttonI.GetName())

		command := buttonI.GetCurrentCommand()

		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					keyboard_names.RENAME_BUTTON,
					callback_data.QueryDataType{
						Keyboard: callback_data.KeyboardType_BUTTONS,
						Path:     callbackDataPath,
						Action:   callback_data.ActionType_EDIT,
					}.GetBase64ProtoString()),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					"Topic: "+command.Topic,
					callback_data.QueryDataType{
						Keyboard: callback_data.KeyboardType_BUTTONS,
						Path:     callbackDataPath,
						Action:   callback_data.ActionType_EDIT_COMMAND_TOPIC,
					}.GetBase64ProtoString()),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					"Value: "+command.Value,
					callback_data.QueryDataType{
						Keyboard: callback_data.KeyboardType_BUTTONS,
						Path:     callbackDataPath,
						Action:   callback_data.ActionType_EDIT_COMMAND_VALUE,
					}.GetBase64ProtoString()),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf(keyboard_names.QOS, command.Qos),
					callback_data.QueryDataType{
						Keyboard: callback_data.KeyboardType_BUTTONS,
						Path:     callbackDataPath,
						Action:   callback_data.ActionType_SWITCH_QOS,
						IntValue: int32(command.Qos+1) % 3,
					}.GetBase64ProtoString()),
				tgbotapi.NewInlineKeyboardButtonData(
					fmt.Sprintf(keyboard_names.RETAINED, command.Retained),
					callback_data.QueryDataType{
						Keyboard:  callback_data.KeyboardType_BUTTONS,
						Path:      callbackDataPath,
						Action:    callback_data.ActionType_SWITCH_RETAINED,
						BoolValue: !command.Retained,
					}.GetBase64ProtoString()),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					keyboard_names.DELETE_BUTTON,
					callback_data.QueryDataType{
						Keyboard: callback_data.KeyboardType_BUTTONS,
						Path:     callbackDataPath,
						Action:   callback_data.ActionType_DELETE,
					}.GetBase64ProtoString()),
			),
		)

	case button_types.TOGGLE:
		text = fmt.Sprintf("Edit %v <code>%v</code>:", buttonI.GetType(), buttonI.GetFullName())

		for commandIndex, command := range buttonI.GetCommands() {
			inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard,
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(
						fmt.Sprintf("Command %v: %v", commandIndex, command.Name),
						callback_data.QueryDataType{
							Keyboard: callback_data.KeyboardType_BUTTONS,
							Path:     callbackDataPath,
							Action:   callback_data.ActionType_EDIT_COMMAND,
							Index:    int32(commandIndex),
						}.GetBase64ProtoString()),
				),
			)
		}

		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					keyboard_names.DELETE_BUTTON,
					callback_data.QueryDataType{
						Keyboard: callback_data.KeyboardType_BUTTONS,
						Path:     callbackDataPath,
						Action:   callback_data.ActionType_DELETE,
					}.GetBase64ProtoString()),
			),
		)
	case button_types.MULTI_VALUE:
		text = fmt.Sprintf("Edit %v <code>%v</code>:", buttonI.GetType(), buttonI.GetName())

		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					keyboard_names.RENAME_BUTTON,
					callback_data.QueryDataType{
						Keyboard: callback_data.KeyboardType_BUTTONS,
						Path:     callbackDataPath,
						Action:   callback_data.ActionType_EDIT,
					}.GetBase64ProtoString()),
			),
		)

		for commandIndex, command := range buttonI.GetCommands() {
			inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard,
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(
						fmt.Sprintf("Command %v: %v", commandIndex, command.Name),
						callback_data.QueryDataType{
							Keyboard: callback_data.KeyboardType_BUTTONS,
							Path:     callbackDataPath,
							Action:   callback_data.ActionType_EDIT_COMMAND,
							Index:    int32(commandIndex),
						}.GetBase64ProtoString()),
				),
			)
		}

		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					keyboard_names.ADD_COMMAND,
					callback_data.QueryDataType{
						Keyboard: callback_data.KeyboardType_BUTTONS,
						Path:     callbackDataPath,
						Action:   callback_data.ActionType_ADD_COMMAND,
					}.GetBase64ProtoString()),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					keyboard_names.DELETE_BUTTON,
					callback_data.QueryDataType{
						Keyboard: callback_data.KeyboardType_BUTTONS,
						Path:     callbackDataPath,
						Action:   callback_data.ActionType_DELETE,
					}.GetBase64ProtoString()),
			),
		)
	case button_types.PRINT_LAST_SUB_VALUE:
		text = fmt.Sprintf("Edit %v <code>%v</code>:", buttonI.GetType(), buttonI.GetName())

		var subscriptionID = buttonI.GetSubscriptions()[0]
		var subscriptionTopic string
		if subscriptionID >= 0 && subscriptionID < len(userSubscriptions) {
			subscriptionTopic = userSubscriptions[subscriptionID].Topic
		} else {
			subscriptionTopic = "(not selected)"
		}

		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					keyboard_names.RENAME_BUTTON,
					callback_data.QueryDataType{
						Keyboard: callback_data.KeyboardType_BUTTONS,
						Path:     callbackDataPath,
						Action:   callback_data.ActionType_EDIT,
					}.GetBase64ProtoString()),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					"Subscription: "+subscriptionTopic,
					callback_data.QueryDataType{
						Keyboard: callback_data.KeyboardType_BUTTONS,
						Path:     callbackDataPath,
						Action:   callback_data.ActionType_SWITCH_SUBSCRIPTION,
						IntValue: int32((subscriptionID + 1) % len(userSubscriptions)),
					}.GetBase64ProtoString()),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					keyboard_names.DELETE_BUTTON,
					callback_data.QueryDataType{
						Keyboard: callback_data.KeyboardType_BUTTONS,
						Path:     callbackDataPath,
						Action:   callback_data.ActionType_DELETE,
					}.GetBase64ProtoString()),
			),
		)
	case button_types.DRAW_CHART:
		text = fmt.Sprintf("Edit %v <code>%v</code>:", buttonI.GetType(), buttonI.GetName())

		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					keyboard_names.RENAME_BUTTON,
					callback_data.QueryDataType{
						Keyboard: callback_data.KeyboardType_BUTTONS,
						Path:     callbackDataPath,
						Action:   callback_data.ActionType_EDIT,
					}.GetBase64ProtoString()),
			),
		)

		for subscriptionButtonIndex, subscriptionUserIndex := range buttonI.GetSubscriptions() {
			var subscription *models.Subscription
			if subscriptionUserIndex >= len(userSubscriptions) {
				subscription = &models.Subscription{
					Topic: "(not found)",
				}
			} else {
				subscription = userSubscriptions[subscriptionUserIndex]
			}
			inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard,
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(
						fmt.Sprintf("Sub chart %v: %v", subscriptionButtonIndex, subscription.Topic),
						callback_data.QueryDataType{
							Keyboard: callback_data.KeyboardType_BUTTONS,
							Path:     callbackDataPath,
							Action:   callback_data.ActionType_SWITCH_SUBSCRIPTION,
							IntValue: int32((subscriptionUserIndex + 1) % len(userSubscriptions)),
							Index:    int32(subscriptionButtonIndex),
						}.GetBase64ProtoString()),
					tgbotapi.NewInlineKeyboardButtonData(
						keyboard_names.DELETE,
						callback_data.QueryDataType{
							Keyboard: callback_data.KeyboardType_BUTTONS,
							Path:     callbackDataPath,
							Action:   callback_data.ActionType_DELETE_SUB_CHART,
							Index:    int32(subscriptionButtonIndex),
						}.GetBase64ProtoString()),
				),
			)
		}

		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					keyboard_names.ADD_SUB_CHART,
					callback_data.QueryDataType{
						Keyboard: callback_data.KeyboardType_BUTTONS,
						Path:     callbackDataPath,
						Action:   callback_data.ActionType_ADD_SUB_CHART,
					}.GetBase64ProtoString()),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					keyboard_names.DELETE_BUTTON,
					callback_data.QueryDataType{
						Keyboard: callback_data.KeyboardType_BUTTONS,
						Path:     callbackDataPath,
						Action:   callback_data.ActionType_DELETE,
					}.GetBase64ProtoString()),
			),
		)
	}

	if len(callbackDataPath) > 0 {
		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard,
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					keyboard_names.BACK,
					callback_data.QueryDataType{
						Keyboard: callback_data.KeyboardType_BUTTONS,
						Path:     callbackDataPath[:len(callbackDataPath)-1],
					}.GetBase64ProtoString()),
			),
		)
	}

	return text, &inlineKeyboard
}

func GetAddButtonKeyboard(buttonType button_types.ButtonType, callbackDataPath []int32) (string, *tgbotapi.InlineKeyboardMarkup) {
	text := "Select the button type below and next send me the button name (you can use emojis)"
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				buttonType.TypeString(),
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_BUTTONS,
					Path:     callbackDataPath,
					Action:   callback_data.ActionType_SWITCH_BUTTON_TYPE,
					IntValue: int32(buttonType.NextType(len(callbackDataPath) > 0)), // skip folder if not in root
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.BACK_TO_LIST,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_BUTTONS,
					Path:     callbackDataPath,
					Action:   callback_data.ActionType_BACK_TO_LIST,
				}.GetBase64ProtoString()),
		),
	)
	return text, &inlineKeyboard
}

func GetEditButtonNameKeyboard(currButton button_interface.ButtonI, callbackDataPath []int32) (string, *tgbotapi.InlineKeyboardMarkup) {
	text := fmt.Sprintf("Current <code>%v</code> name is <code>%v</code>", currButton.GetType().String(), currButton.GetName())
	text += "\nSend me the new name if you want to rename it (you can use emojis)"

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.BACK_TO_MENU,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_BUTTONS,
					Path:     callbackDataPath,
					Action:   callback_data.ActionType_BACK_TO_MENU,
				}.GetBase64ProtoString()),
		),
	)
	return text, &inlineKeyboard
}

func GetDeleteButtonKeyboard(currButton button_interface.ButtonI, callbackDataPath []int32) (string, *tgbotapi.InlineKeyboardMarkup) {
	text := fmt.Sprintf("Are you sure you want to delete %v <code>%v</code>?", currButton.GetType(), currButton.GetName())
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.YES,
				callback_data.QueryDataType{
					Keyboard:  callback_data.KeyboardType_BUTTONS,
					Path:      callbackDataPath,
					Action:    callback_data.ActionType_DELETE,
					BoolValue: true,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.BACK,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_BUTTONS,
					Path:     callbackDataPath,
				}.GetBase64ProtoString()),
		),
	)
	return text, &inlineKeyboard
}

func GetMultiValueCommandsKeyboard(currentButton button_interface.ButtonI) (string, *tgbotapi.InlineKeyboardMarkup) {
	inlineText := fmt.Sprintf("Select command for %v:", currentButton.GetName())

	if len(currentButton.GetCommands()) > 0 {
		currentPath := make([]int32, 0)

		parent := currentButton.GetParent()
		for parent != nil {
			for index, button := range *parent.GetButtons() {
				if button == currentButton {
					currentPath = append(currentPath, int32(index))
					break
				}
			}
			parent = parent.GetParent()
		}

		// reverse slice
		for left, right := 0, len(currentPath)-1; left < right; left, right = left+1, right-1 {
			currentPath[left], currentPath[right] = currentPath[right], currentPath[left]
		}

		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup()
		for commandIndex, command := range currentButton.GetCommands() {
			inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard,
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(
						command.Name,
						callback_data.QueryDataType{
							Keyboard: callback_data.KeyboardType_COMMAND,
							Path:     currentPath,
							Index:    int32(commandIndex),
						}.GetBase64ProtoString()),
				),
			)
		}
		return inlineText, &inlineKeyboard
	}

	inlineText += fmt.Sprintf("\n(no commands, you need to configure button in %v -> %v)", button_names.SETTINGS, button_names.EDIT_BUTTONS)

	return inlineText, nil
}

func GetCommandAddKeyboard(currButton button_interface.ButtonI, callbackDataPath []int32, qos byte, retained bool) (string, *tgbotapi.InlineKeyboardMarkup) {
	text := "Select QoS level and Retained flag of the new command"
	text += fmt.Sprintf("\nNext send me the name to add command into %v <code>%v</code>", currButton.GetType(), currButton.GetName())

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(keyboard_names.QOS, qos),
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_BUTTONS,
					Path:     callbackDataPath,
					Action:   callback_data.ActionType_SWITCH_QOS,
					IntValue: int32(qos+1) % 3,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(keyboard_names.RETAINED, retained),
				callback_data.QueryDataType{
					Keyboard:  callback_data.KeyboardType_BUTTONS,
					Path:      callbackDataPath,
					Action:    callback_data.ActionType_SWITCH_RETAINED,
					BoolValue: !retained,
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.BACK,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_BUTTONS,
					Path:     callbackDataPath,
				}.GetBase64ProtoString()),
		),
	)
	return text, &inlineKeyboard
}

func GetCommandEditKeyboard(currButton button_interface.ButtonI, commandIndex int32, callbackDataPath []int32) (string, *tgbotapi.InlineKeyboardMarkup) {
	if currButton.GetType() == button_types.SINGLE_VALUE {
		return GetButtonsKeyboard(currButton, callbackDataPath, nil)
	}

	text := fmt.Sprintf("Edit command of %v <code>%v</code>", currButton.GetType(), currButton.GetFullName())

	command := currButton.GetCommands()[commandIndex]

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Name: "+command.Name,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_BUTTONS,
					Path:     callbackDataPath,
					Action:   callback_data.ActionType_EDIT_COMMAND_NAME,
					Index:    commandIndex,
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Topic: "+command.Topic,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_BUTTONS,
					Path:     callbackDataPath,
					Action:   callback_data.ActionType_EDIT_COMMAND_TOPIC,
					Index:    commandIndex,
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"Value: "+command.Value,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_BUTTONS,
					Path:     callbackDataPath,
					Action:   callback_data.ActionType_EDIT_COMMAND_VALUE,
					Index:    commandIndex,
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(keyboard_names.QOS, command.Qos),
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_BUTTONS,
					Path:     callbackDataPath,
					Action:   callback_data.ActionType_SWITCH_QOS,
					IntValue: int32(command.Qos+1) % 3,
					Index:    commandIndex,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(keyboard_names.RETAINED, command.Retained),
				callback_data.QueryDataType{
					Keyboard:  callback_data.KeyboardType_BUTTONS,
					Path:      callbackDataPath,
					Action:    callback_data.ActionType_SWITCH_RETAINED,
					BoolValue: !command.Retained,
					Index:     commandIndex,
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.DELETE_COMMAND,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_BUTTONS,
					Path:     callbackDataPath,
					Action:   callback_data.ActionType_DELETE_COMMAND,
					Index:    commandIndex,
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.BACK,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_BUTTONS,
					Path:     callbackDataPath,
				}.GetBase64ProtoString()),
		),
	)
	return text, &inlineKeyboard
}

func GetDeleteCommandKeyboard(currButton button_interface.ButtonI, commandIndex int32, callbackDataPath []int32) (string, *tgbotapi.InlineKeyboardMarkup) {
	command := currButton.GetCommands()[commandIndex]
	text := fmt.Sprintf("Are you sure you want to delete the command <code>%v</code>?", command.Name)
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.YES,
				callback_data.QueryDataType{
					Keyboard:  callback_data.KeyboardType_BUTTONS,
					Path:      callbackDataPath,
					Action:    callback_data.ActionType_DELETE_COMMAND,
					BoolValue: true,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.BACK,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_BUTTONS,
					Path:     callbackDataPath,
					Action:   callback_data.ActionType_EDIT_COMMAND,
					Index:    commandIndex,
				}.GetBase64ProtoString()),
		),
	)
	return text, &inlineKeyboard
}

func GetDeleteSubscriptionKeyboard(currButton button_interface.ButtonI, subscriptionButtonIndex int32, userSubscriptions []*models.Subscription, callbackDataPath []int32) (string, *tgbotapi.InlineKeyboardMarkup) {
	if int(subscriptionButtonIndex) >= len(currButton.GetSubscriptions()) {
		return "", nil
	}
	subscriptionUserIndex := currButton.GetSubscriptions()[subscriptionButtonIndex]
	if subscriptionUserIndex >= len(userSubscriptions) {
		return "", nil
	}
	subscription := userSubscriptions[subscriptionUserIndex]

	text := fmt.Sprintf("Are you sure you want to delete the subscription <code>%v</code> from chart?", subscription.Topic)
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.YES,
				callback_data.QueryDataType{
					Keyboard:  callback_data.KeyboardType_BUTTONS,
					Path:      callbackDataPath,
					Action:    callback_data.ActionType_DELETE_SUB_CHART,
					BoolValue: true,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.BACK,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_BUTTONS,
					Path:     callbackDataPath,
				}.GetBase64ProtoString()),
		),
	)
	return text, &inlineKeyboard
}

func GetEditCommandNameKeyboard(currButton button_interface.ButtonI, commandIndex int32, callbackDataPath []int32) (string, *tgbotapi.InlineKeyboardMarkup) {
	command := currButton.GetCommands()[commandIndex]
	text := fmt.Sprintf("The current command name is <code>%v</code>", command.Name)
	text += "\nSend me the new name if you want to rename it"

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.BACK,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_BUTTONS,
					Path:     callbackDataPath,
					Action:   callback_data.ActionType_EDIT_COMMAND,
					Index:    commandIndex,
				}.GetBase64ProtoString()),
		),
	)
	return text, &inlineKeyboard
}

func GetEditCommandTopicKeyboard(currButton button_interface.ButtonI, commandIndex int32, callbackDataPath []int32) (string, *tgbotapi.InlineKeyboardMarkup) {
	command := currButton.GetCommands()[commandIndex]
	topic := command.Topic
	if len(topic) == 0 {
		topic = "empty"
	}

	text := fmt.Sprintf("The current <code>%v</code> command topic is <code>%v</code>", command.Name, topic)
	text += "\nSend me the new topic if you want to change it"

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.BACK,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_BUTTONS,
					Path:     callbackDataPath,
					Action:   callback_data.ActionType_EDIT_COMMAND,
					Index:    commandIndex,
				}.GetBase64ProtoString()),
		),
	)
	return text, &inlineKeyboard
}
func GetEditCommandValueKeyboard(currButton button_interface.ButtonI, commandIndex int32, callbackDataPath []int32) (string, *tgbotapi.InlineKeyboardMarkup) {
	command := currButton.GetCommands()[commandIndex]
	topic := command.Topic
	if len(topic) == 0 {
		topic = "empty"
	}
	value := command.Value
	if len(value) == 0 {
		value = "empty"
	}

	text := fmt.Sprintf("Current <code>%v</code> command <code>%v</code> topic value is <code>%v</code>", command.Name, topic, value)
	text += "\nSend me the new value if you want to change it"

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.BACK,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_BUTTONS,
					Path:     callbackDataPath,
					Action:   callback_data.ActionType_EDIT_COMMAND,
					Index:    commandIndex,
				}.GetBase64ProtoString()),
		),
	)
	return text, &inlineKeyboard
}
