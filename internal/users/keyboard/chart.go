package keyboard

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"mqtg-bot/internal/users/keyboard/callback_data"
	"mqtg-bot/internal/users/menu/button_interface"
)

func GetShowChartKeyboard(currButton button_interface.ButtonI, callbackDataPath []int32) (string, *tgbotapi.InlineKeyboardMarkup) {
	text := fmt.Sprintf("Select a range for the chart <code>%v</code>", currButton.GetName())
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"30 min",
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_CHART,
					Path:     callbackDataPath,
					IntValue: 30,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				"2 hours",
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_CHART,
					Path:     callbackDataPath,
					IntValue: 120,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				"6 hours",
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_CHART,
					Path:     callbackDataPath,
					IntValue: 360,
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"1 day",
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_CHART,
					Path:     callbackDataPath,
					IntValue: 1440,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				"1 week",
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_CHART,
					Path:     callbackDataPath,
					IntValue: 10080,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				"all",
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_CHART,
					Path:     callbackDataPath,
					IntValue: 0,
				}.GetBase64ProtoString()),
		),
	)
	return text, &inlineKeyboard
}
