package internal

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"mqtg-bot/internal/common"
	"net/http"
)

func (bot *TelegramBot) SendAnswer(chatID int64, answer *common.BotMessage) {
	if answer == nil {
		return
	}

	if answer.MessageID != 0 { // need edit existing message
		if len(answer.InlineText) > 0 {
			bot.EditMessageText(chatID, answer.MessageID, answer.InlineText)
		}

		newInlineKeyboard, ok := answer.InlineKeyboard.(*tgbotapi.InlineKeyboardMarkup)
		if ok && newInlineKeyboard != nil {
			bot.EditInlineKeyboard(chatID, answer.MessageID, newInlineKeyboard)
		}
	} else {
		if len(answer.Photo) > 0 {
			bot.NewPhotoUpload(chatID, answer.MainText, answer.Photo, answer.InlineKeyboard)
		} else {
			if answer.MainText != "" {
				bot.SendMessage(chatID, answer.MainText, answer.MainMenu)
			}
			if answer.InlineText != "" {
				bot.SendMessage(chatID, answer.InlineText, answer.InlineKeyboard)
			}
		}
	}
}

func (bot *TelegramBot) NewPhotoUpload(chatID int64, text string, payload []byte, replyMarkup interface{}) {
	fileBytes := tgbotapi.FileBytes{
		Bytes: payload,
	}
	photoConfig := tgbotapi.NewPhotoUpload(chatID, fileBytes)
	photoConfig.Caption = text
	photoConfig.ReplyMarkup = replyMarkup
	photoConfig.ParseMode = tgbotapi.ModeHTML
	bot.ConfigureAndSend(photoConfig)
}

func (bot *TelegramBot) EditMessageText(chatID int64, messageID int, text string) {
	msg := tgbotapi.NewEditMessageText(
		chatID,
		messageID,
		text,
	)
	bot.ConfigureAndSend(msg)
}

func (bot *TelegramBot) EditPhotoMessage(chatID int64, messageID int, photo []byte) {
	msg := tgbotapi.EditMessageMediaConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    chatID,
			MessageID: messageID,
		},
		Media: tgbotapi.NewInputMediaPhoto("https://upload.wikimedia.org/wikipedia/ru/3/31/Winniethepooh.jpg"),
	}
	bot.ConfigureAndSend(msg)
}

func (bot *TelegramBot) EditInlineKeyboard(chatID int64, messageID int, newInlineKeyboard *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewEditMessageReplyMarkup(
		chatID,
		messageID,
		*newInlineKeyboard,
	)
	bot.ConfigureAndSend(msg)
}

func (bot *TelegramBot) SendMessage(chatID int64, text string, replyMarkup interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	inlineKeyboard, ok := replyMarkup.(*tgbotapi.InlineKeyboardMarkup)
	if ok {
		if inlineKeyboard != nil && len(inlineKeyboard.InlineKeyboard) > 0 {
			msg.ReplyMarkup = inlineKeyboard
		}
	} else {
		msg.ReplyMarkup = replyMarkup
	}
	bot.ConfigureAndSend(msg)
}

func (bot *TelegramBot) DownloadPhoto(photo []tgbotapi.PhotoSize) ([]byte, error) {
	link, err := bot.GetFileDirectURL(photo[len(photo)-1].FileID)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (bot *TelegramBot) ConfigureAndSend(msg tgbotapi.Chattable) {
	emtc, ok := msg.(tgbotapi.EditMessageTextConfig)
	if ok {
		emtc.ParseMode = tgbotapi.ModeHTML
		bot.Send(emtc)
		return
	}

	mc, ok := msg.(tgbotapi.MessageConfig)
	if ok {
		mc.ParseMode = tgbotapi.ModeHTML
		bot.Send(mc)
		return
	}

	bot.Send(msg)
}

func (bot *TelegramBot) Send(msg tgbotapi.Chattable) {
	bot.metrics.numOfOutMessagesToTelegram.Inc()

	_, err := bot.BotAPI.Send(msg)
	if err != nil {
		log.Printf("Send error: %v. Msg: %+v", err, msg)
	}
}
