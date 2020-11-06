package menu

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"mqtg-bot/internal/users/menu/button_interface"
	"mqtg-bot/internal/users/menu/button_names"
	"mqtg-bot/internal/users/menu/button_types"
)

func (menu *MainMenu) Back() {
	switch menu.CurrPath.GetParent() {
	case nil, &menu.UserButtons, &menu.CommonButtons:
		menu.ResetCurrentPath()
	default:
		menu.CurrPath = menu.CurrPath.GetParent()
	}
}

func (menu *MainMenu) SetPressedButtonLikeCurrentPath(message string, currPath button_interface.ButtonI) button_interface.ButtonI {
	var button button_interface.ButtonI
	if currPath == &menu.CommonButtons { // if in main menu
		button = menu.findInButtons(menu.UserButtons.GetButtons(), message) // try to find in UserButtons first
		if button == nil {
			button = menu.findInButtons(menu.CommonButtons.GetButtons(), message) // and in CommonButtons next
		}
	} else {
		button = menu.findInButtons(currPath.GetButtons(), message) // else in currPath
	}

	if button == nil {
		menu.Back()
	}

	return button
}

func (menu *MainMenu) findInButtons(buttons *[]button_interface.ButtonI, message string) button_interface.ButtonI {
	if buttons != nil {
		for _, button := range *buttons {
			if button.GetName() == message {
				if button.GetType() == button_types.FOLDER {
					menu.CurrPath = button
				}
				return button
			}
		}
	}

	return nil
}

func (menu *MainMenu) GetCurrPathMainMenu() *tgbotapi.ReplyKeyboardMarkup {
	var keyboardButtons []tgbotapi.KeyboardButton

	if menu.CurrPath == &menu.CommonButtons {
		for _, button := range *menu.UserButtons.GetButtons() {
			keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButton(button.GetName()))
		}

		for _, button := range *menu.CommonButtons.GetButtons() {
			keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButton(button.GetName()))
		}
	} else {
		for _, button := range *menu.CurrPath.GetButtons() {
			keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButton(button.GetName()))
		}

		keyboardButtons = append(keyboardButtons, tgbotapi.NewKeyboardButton(button_names.BACK))
	}

	return splitAndConvertIntoMainMenuKeyboardMarkup(keyboardButtons)
}

func splitAndConvertIntoMainMenuKeyboardMarkup(keyboardButtons []tgbotapi.KeyboardButton) *tgbotapi.ReplyKeyboardMarkup {
	buttonsCount := len(keyboardButtons)

	var numRows = 3
	if buttonsCount > 9 {
		numRows = 4
	}

	splitBy := (buttonsCount-1)/numRows + 1
	rowsCount := (buttonsCount-1)/splitBy + 1

	keyboardButtonRows := make([][]tgbotapi.KeyboardButton, 0, rowsCount)
	for index, keyboardButton := range keyboardButtons {
		row := index / splitBy
		if len(keyboardButtonRows) <= row {
			keyboardButtonRows = append(keyboardButtonRows, make([]tgbotapi.KeyboardButton, 0, splitBy))
		}
		keyboardButtonRows[row] = append(keyboardButtonRows[row], keyboardButton)
	}

	keyboard := tgbotapi.NewReplyKeyboard(keyboardButtonRows...)

	return &keyboard
}
