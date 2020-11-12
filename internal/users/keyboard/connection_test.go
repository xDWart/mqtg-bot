package keyboard

import (
	"encoding/base64"
	"github.com/golang/protobuf/proto"
	"mqtg-bot/internal/users/keyboard/callback_data"
	"testing"
)

func TestGetConnectionStringKeyboard(t *testing.T) {
	testMqttUrl := "tcp://user:password@host:port/path"
	inlineText, inlineKeyboard := GetConnectionStringKeyboard(testMqttUrl)

	switch true {
	case len(inlineText) == 0:
		t.Errorf("unexpected empty inline text")
	case inlineKeyboard == nil:
		t.Errorf("unexpected nil keyboard")
	case len(inlineKeyboard.InlineKeyboard) == 0,
		len(inlineKeyboard.InlineKeyboard[0]) == 0:
		t.Errorf("unexpected keyboard len")
	}

	if inlineKeyboard.InlineKeyboard[0][0].Text != testMqttUrl {
		t.Errorf("unexpected keyboard text: %v", inlineKeyboard.InlineKeyboard[0][0].Text)
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(*inlineKeyboard.InlineKeyboard[0][0].CallbackData)
	if err != nil {
		t.Errorf("base64 decode error: %v", err)
	}

	var callbackData callback_data.QueryDataType
	err = proto.Unmarshal(decodedBytes, &callbackData)
	if err != nil {
		t.Errorf("proto unmarshall error: %v", err)
	}

	if callbackData.Keyboard != callback_data.KeyboardType_CONNECTION {
		t.Errorf("unexpected keyboard type: %v", callbackData.Keyboard)
	}

	inlineText, inlineKeyboard = GetConnectionStringKeyboard("")
	switch {
	case len(inlineText) == 0:
		t.Errorf("unexpected empty inline text")
	case inlineKeyboard != nil:
		t.Errorf("inline keyboard not nil: %#v", inlineKeyboard)
	}
}
