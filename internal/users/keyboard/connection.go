package keyboard

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"mqtg-bot/internal/users/keyboard/callback_data"
)

func GetConnectionStringKeyboard(mqttUrl string) (string, *tgbotapi.InlineKeyboardMarkup) {
	inlineText := "Send me your MQTT broker connection URL in the following format of <code>(tcp|ssl|ws|wss)://user:password@host:port/path</code>"

	if len(mqttUrl) > 0 {
		inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(
					mqttUrl,
					callback_data.QueryDataType{
						Keyboard: callback_data.KeyboardType_CONNECTION,
					}.GetBase64ProtoString()),
			),
		)
		return inlineText, &inlineKeyboard
	}

	return inlineText, nil
}
