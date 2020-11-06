package keyboard

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"mqtg-bot/internal/users/keyboard/callback_data"
	"mqtg-bot/internal/users/keyboard/keyboard_names"
)

func GetPublishKeyboard(qos byte, retained bool) (string, *tgbotapi.InlineKeyboardMarkup) {
	inlineText := "Select QoS level and Retained flag"
	inlineText += "\nand send me the <code>topic</code> for publishing"

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(keyboard_names.QOS, qos),
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_PUBLISH,
					Action:   callback_data.ActionType_SWITCH_QOS,
					IntValue: int32(qos+1) % 3,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(keyboard_names.RETAINED, retained),
				callback_data.QueryDataType{
					Keyboard:  callback_data.KeyboardType_PUBLISH,
					Action:    callback_data.ActionType_SWITCH_RETAINED,
					BoolValue: !retained,
				}.GetBase64ProtoString()),
		),
	)

	return inlineText, &inlineKeyboard
}
