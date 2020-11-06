package menu

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"mqtg-bot/internal/users/menu/button_interface"
	"mqtg-bot/internal/users/menu/button_names"
)

var ConfigureConnectionMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(button_names.CONFIGURE_CONNECTION),
	),
)

var commonMenu = button_interface.FolderButton{
	Name: button_names.MAIN_MENU,
	Buttons: []button_interface.ButtonI{
		&button_interface.FolderButton{
			Name: button_names.SETTINGS,
			Buttons: []button_interface.ButtonI{
				&button_interface.SystemButton{
					Name: button_names.PUBLISH,
				},
				&button_interface.SystemButton{
					Name: button_names.SUBSCRIPTIONS,
				},
				&button_interface.SystemButton{
					Name: button_names.EDIT_BUTTONS,
				},
				&button_interface.SystemButton{
					Name: button_names.DISCONNECT,
				},
			},
		},
	},
}
