package common

type BotMessage struct {
	MessageID int // for edit existing message

	MainText string
	MainMenu interface{}

	InlineText     string
	InlineKeyboard interface{}

	Photo []byte
}
