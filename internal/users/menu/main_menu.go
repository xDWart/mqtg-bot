package menu

import (
	"encoding/json"
	"log"
	"mqtg-bot/internal/users/menu/button_interface"
)

type MainMenu struct {
	UserButtons   button_interface.FolderButton
	CommonButtons button_interface.FolderButton
	CurrPath      button_interface.ButtonI
}

func (menu *MainMenu) ResetCurrentPath() {
	menu.CurrPath = &menu.CommonButtons
}

func (menu *MainMenu) AppendCommonMenuAndSetParentLinks() {
	menu.CommonButtons = commonMenu
	menu.UserButtons.SetParent(nil)
	menu.CommonButtons.SetParent(nil)
	menu.ResetCurrentPath()
}

func (menu *MainMenu) GenerateJsonb() ([]byte, error) {
	return json.Marshal(menu.UserButtons)
}

func (menu *MainMenu) LoadMenuFromJsonb(data []byte) {
	err := json.Unmarshal(data, &menu.UserButtons)
	if err != nil {
		log.Printf("Unmarshal error: %v", err)
	}
}
