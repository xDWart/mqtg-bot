package keyboard

import (
	"testing"
)

func TestGetPublishKeyboard(t *testing.T) {
	var qos byte
	var retained bool
	inlineText, inlineKeyboard := GetPublishKeyboard(qos, retained)

	switch true {
	case len(inlineText) == 0:
		t.Errorf("unexpected empty inline text")
	case inlineKeyboard == nil:
		t.Errorf("unexpected nil keyboard")
	case len(inlineKeyboard.InlineKeyboard) == 0,
		len(inlineKeyboard.InlineKeyboard[0]) == 0:
		t.Errorf("unexpected keyboard len")
	}
}
