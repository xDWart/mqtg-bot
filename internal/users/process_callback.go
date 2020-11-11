package users

import (
	"bytes"
	"encoding/base64"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/golang/protobuf/proto"
	"github.com/wcharczuk/go-chart"
	"log"
	"mqtg-bot/internal/common"
	"mqtg-bot/internal/models"
	"mqtg-bot/internal/users/keyboard"
	"mqtg-bot/internal/users/keyboard/callback_data"
	"mqtg-bot/internal/users/menu/button_interface"
	"mqtg-bot/internal/users/menu/button_types"
	"mqtg-bot/internal/users/state"
	"strconv"
	"strings"
	"time"
)

func (user *User) ProcessCallback(callbackQuery *tgbotapi.CallbackQuery) *common.BotMessage {
	var callbackData callback_data.QueryDataType
	var newText string
	var newInlineKeyboard *tgbotapi.InlineKeyboardMarkup

	// NB: iOS Telegram repeats callbacks on device blocking
	if user.lastCallbackData == callbackQuery.Data &&
		time.Now().Sub(user.lastCallbackTime) < 2*time.Second {
		// simple defense from multi callback
		return nil
	}
	user.lastCallbackData = callbackQuery.Data
	user.lastCallbackTime = time.Now()

	decodedBytes, err := base64.StdEncoding.DecodeString(callbackQuery.Data)
	if err == nil {
		err = proto.Unmarshal(decodedBytes, &callbackData)
	}
	callbackData.MessageId = int64(callbackQuery.Message.MessageID)

	log.Printf("Telegram received callback from %v (Chat.ID %v): %+v",
		callbackQuery.From.UserName, callbackQuery.Message.Chat.ID, callbackData)

	switch callbackData.Keyboard {
	case callback_data.KeyboardType_CONNECTION:
		return user.processConnectionString("")

	case callback_data.KeyboardType_SUBSCRIPTIONS:
		newText, newInlineKeyboard, err = user.processSubscriptionCallback(&callbackData)

	case callback_data.KeyboardType_PUBLISH:
		newInlineKeyboard, err = user.processPublishCallback(&callbackData)

	case callback_data.KeyboardType_BUTTONS:
		newText, newInlineKeyboard = user.processButtonMenuCallback(&callbackData)

	case callback_data.KeyboardType_CHART:
		mainText, photo := user.processChartMenuCallback(&callbackData)
		return &common.BotMessage{
			MainText: mainText,
			Photo:    photo,
		}

	case callback_data.KeyboardType_COMMAND:
		var currButton button_interface.ButtonI
		currButton = &user.menu.UserButtons
		for _, buttonIndex := range callbackData.Path {
			if int(buttonIndex) < len(*currButton.GetButtons()) {
				currButton = (*currButton.GetButtons())[buttonIndex]
			}
		}
		if int(callbackData.Index) < len(currButton.GetCommands()) {
			command := currButton.GetCommands()[callbackData.Index]
			text := fmt.Sprintf("%v (%v)", currButton.GetName(), command.Name)
			if len(command.Topic) > 0 && len(command.Value) > 0 {
				user.mqtt.Publish(command.Topic, command.Qos, command.Retained, command.Value)
			} else {
				text += " - nothing to send, you need to configure this command"
			}
			return &common.BotMessage{
				MainText: text,
			}
		}
	}

	if err != nil || (len(newText) == 0 && newInlineKeyboard == nil) {
		mainText := fmt.Sprintf("Bad CallbackQuery.Data: %+v", callbackData)
		if err != nil {
			mainText += fmt.Sprintf(". Error: %v", err)
		}
		return &common.BotMessage{
			MainText: mainText,
			MainMenu: user.getMainMenu(),
		}
	}

	return &common.BotMessage{
		MessageID:      callbackQuery.Message.MessageID,
		InlineText:     newText,
		InlineKeyboard: newInlineKeyboard,
	}
}

func (user *User) processSubscriptionCallback(callbackData *callback_data.QueryDataType) (string, *tgbotapi.InlineKeyboardMarkup, error) {
	switch callbackData.Action {
	case callback_data.ActionType_ADD_SUBSCRIPTION:
		newText, newInlineKeyboard := keyboard.GetAddSubscriptionKeyboard(user.state.SubscriptionType, user.state.Qos, user.state.SubscriptionDataType)
		user.state = state.StateStruct{
			State: state.ADD_SUBSCRIPTION_STATE,
		}
		return newText, newInlineKeyboard, nil

	case callback_data.ActionType_SWITCH_ADDED_SUBSCRIPTION_TYPE:
		user.state.SubscriptionType = models.SubscriptionType(callbackData.IntValue)
		_, newInlineKeyboard := keyboard.GetAddSubscriptionKeyboard(user.state.SubscriptionType, user.state.Qos, user.state.SubscriptionDataType)
		return "", newInlineKeyboard, nil

	case callback_data.ActionType_SWITCH_ADDED_SUBSCRIPTION_QOS:
		user.state.Qos = byte(callbackData.IntValue)
		_, newInlineKeyboard := keyboard.GetAddSubscriptionKeyboard(user.state.SubscriptionType, user.state.Qos, user.state.SubscriptionDataType)
		return "", newInlineKeyboard, nil

	case callback_data.ActionType_SWITCH_ADDED_SUBSCRIPTION_DATA_TYPE:
		user.state.SubscriptionDataType = models.SubscriptionDataType(callbackData.IntValue)
		_, newInlineKeyboard := keyboard.GetAddSubscriptionKeyboard(user.state.SubscriptionType, user.state.Qos, user.state.SubscriptionDataType)
		return "", newInlineKeyboard, nil

	case callback_data.ActionType_BACK_TO_LIST:
		user.state.Reset()
		newText, newInlineKeyboard := keyboard.GetSubscriptionsKeyboard(user.Subscriptions)
		return newText, newInlineKeyboard, nil
	}

	if len(callbackData.Path) < 1 {
		return "", nil, fmt.Errorf("bad subscription index")
	}

	subscriptionIndex := callbackData.Path[0]
	if len(user.Subscriptions) <= int(subscriptionIndex) {
		return "", nil, fmt.Errorf("could not get subscription by index %v", subscriptionIndex)
	}

	subscription := user.Subscriptions[subscriptionIndex]

	switch callbackData.Action {
	case callback_data.ActionType_SWITCH_QOS:
		user.mqtt.Unsubscribe(subscription)
		subscription.Qos = byte(callbackData.IntValue)
		user.mqtt.Subscribe(subscription)
		user.db.Save(&subscription)

	case callback_data.ActionType_SWITCH_SUB_DATA_TYPE:
		user.mqtt.Unsubscribe(subscription)
		subscription.DataType = models.SubscriptionDataType(callbackData.IntValue)
		user.mqtt.Subscribe(subscription)
		user.db.Save(&subscription)

	case callback_data.ActionType_SWITCH_SUBSCRIPTION_TYPE:
		user.mqtt.Unsubscribe(subscription)
		subscription.SubscriptionType = models.SubscriptionType(callbackData.IntValue)
		user.mqtt.Subscribe(subscription)
		user.db.Save(&subscription)

	case callback_data.ActionType_EDIT:
		newText, newInlineKeyboard := keyboard.GetSubscriptionTopicEditKeyboard(subscription, subscriptionIndex)
		user.state = state.StateStruct{
			State:         state.EDIT_SUBSCRIPTION_TOPIC_STATE,
			EditableIndex: subscriptionIndex,
			LastMessageID: int(callbackData.MessageId),
		}
		return newText, newInlineKeyboard, nil

	case callback_data.ActionType_DELETE:
		var newText string
		var newInlineKeyboard *tgbotapi.InlineKeyboardMarkup

		if !callbackData.BoolValue {
			newText, newInlineKeyboard = keyboard.GetSubscriptionDeleteKeyboard(subscription, subscriptionIndex)
		} else { // Yes
			user.unsubscribe(subscription, subscriptionIndex)
			newText, newInlineKeyboard = keyboard.GetSubscriptionsKeyboard(user.Subscriptions)
		}
		return newText, newInlineKeyboard, nil

	case callback_data.ActionType_BEFORE_VALUE_TEXT, callback_data.ActionType_AFTER_VALUE_TEXT:
		newText, newInlineKeyboard := keyboard.GetSubscriptionBeforeAfterValueTextEditKeyboard(callbackData.Action, subscription, subscriptionIndex)

		user.state = state.StateStruct{
			EditableIndex: subscriptionIndex,
			LastMessageID: int(callbackData.MessageId),
		}

		if callbackData.Action == callback_data.ActionType_BEFORE_VALUE_TEXT {
			user.state.State = state.EDIT_BEFORE_VALUE_MESSAGE_TEXT_STATE
		} else {
			user.state.State = state.EDIT_AFTER_VALUE_MESSAGE_TEXT_STATE
		}

		return newText, newInlineKeyboard, nil

	case callback_data.ActionType_EDIT_JSON_PATH:
		newText, newInlineKeyboard := keyboard.GetSubscriptionJsonPathEditKeyboard(subscription, subscriptionIndex)
		user.state = state.StateStruct{
			EditableIndex: subscriptionIndex,
			LastMessageID: int(callbackData.MessageId),
			State:         state.EDIT_JSON_PATH_STATE,
		}
		return newText, newInlineKeyboard, nil

	case callback_data.ActionType_BACK_TO_MENU:
		user.state.Reset()
	}

	newText, newInlineKeyboard := keyboard.GetSubscriptionEditKeyboard(subscription, subscriptionIndex)
	return newText, newInlineKeyboard, nil
}

func (user *User) processPublishCallback(callbackData *callback_data.QueryDataType) (*tgbotapi.InlineKeyboardMarkup, error) {
	switch callbackData.Action {
	case callback_data.ActionType_SWITCH_QOS:
		user.state.Qos = byte(callbackData.IntValue)
	case callback_data.ActionType_SWITCH_RETAINED:
		user.state.Retained = callbackData.BoolValue
	default:
		return nil, fmt.Errorf("unknown action: %v", callbackData)
	}
	_, newInlineKeyboard := keyboard.GetPublishKeyboard(user.state.Qos, user.state.Retained)
	return newInlineKeyboard, nil
}

func (user *User) processButtonMenuCallback(callbackData *callback_data.QueryDataType) (string, *tgbotapi.InlineKeyboardMarkup) {
	var currButton button_interface.ButtonI
	currButton = &user.menu.UserButtons

	for _, buttonIndex := range callbackData.Path {
		if int(buttonIndex) < len(*currButton.GetButtons()) {
			currButton = (*currButton.GetButtons())[buttonIndex]
		}
	}

	switch callbackData.Action {
	case callback_data.ActionType_ADD_BUTTON:
		buttonType := button_types.FOLDER
		skipFolder := len(callbackData.Path) > 0 // skip folder if not in root
		if skipFolder {
			buttonType.NextType(skipFolder)
		}

		user.state = state.StateStruct{
			State:          state.ADD_BUTTON_STATE,
			CurrPath:       callbackData.Path,
			EditableButton: currButton,
			ButtonType:     buttonType,
		}
		return keyboard.GetAddButtonKeyboard(user.state.ButtonType, callbackData.Path)

	case callback_data.ActionType_SWITCH_BUTTON_TYPE:
		user.state.ButtonType = button_types.ButtonType(callbackData.IntValue)
		_, newInlineKeyboard := keyboard.GetAddButtonKeyboard(user.state.ButtonType, callbackData.Path)
		return "", newInlineKeyboard

	case callback_data.ActionType_EDIT: // rename
		user.state = state.StateStruct{
			State:          state.EDIT_BUTTON_NAME_STATE,
			EditableButton: currButton,
			CurrPath:       callbackData.Path,
			LastMessageID:  int(callbackData.MessageId),
		}

		return keyboard.GetEditButtonNameKeyboard(currButton, callbackData.Path)

	case callback_data.ActionType_DELETE:
		if !callbackData.BoolValue {
			return keyboard.GetDeleteButtonKeyboard(currButton, callbackData.Path)
		}

		// Yes
		buttonIndex := callbackData.Path[len(callbackData.Path)-1]
		currButton = currButton.GetParent()
		currButton.DelButton(buttonIndex)
		user.SaveMenuIntoDB()

		callbackData.Path = callbackData.Path[:len(callbackData.Path)-1]

	case callback_data.ActionType_ADD_COMMAND:
		user.state = state.StateStruct{
			State:          state.ADD_NEW_COMMAND,
			EditableButton: currButton,
			CurrPath:       callbackData.Path,
			LastMessageID:  int(callbackData.MessageId),
			Qos:            0,
			Retained:       false,
		}
		return keyboard.GetCommandAddKeyboard(currButton, callbackData.Path, user.state.Qos, user.state.Retained)

	case callback_data.ActionType_SWITCH_QOS, callback_data.ActionType_SWITCH_RETAINED:
		if user.state.State == state.ADD_NEW_COMMAND {
			if callbackData.Action == callback_data.ActionType_SWITCH_QOS {
				user.state.Qos = byte(callbackData.IntValue)
			} else { // callback_data.ActionType_SWITCH_RETAINED
				user.state.Retained = callbackData.BoolValue
			}

			return keyboard.GetCommandAddKeyboard(currButton, callbackData.Path, user.state.Qos, user.state.Retained)
		} else {
			command := currButton.GetCommands()[callbackData.Index]
			if callbackData.Action == callback_data.ActionType_SWITCH_QOS {
				command.Qos = byte(callbackData.IntValue)
			} else { // callback_data.ActionType_SWITCH_RETAINED
				command.Retained = callbackData.BoolValue
			}
			user.SaveMenuIntoDB()

			return keyboard.GetCommandEditKeyboard(currButton, callbackData.Index, callbackData.Path)
		}

	case callback_data.ActionType_EDIT_COMMAND:
		user.state.Reset()
		return keyboard.GetCommandEditKeyboard(currButton, callbackData.Index, callbackData.Path)

	case callback_data.ActionType_DELETE_COMMAND:
		if !callbackData.BoolValue {
			return keyboard.GetDeleteCommandKeyboard(currButton, callbackData.Index, callbackData.Path)
		}

		// Yes
		currButton.DeleteCommand(int(callbackData.Index))
		user.SaveMenuIntoDB()

	case callback_data.ActionType_EDIT_COMMAND_NAME:
		user.state = state.StateStruct{
			State:          state.EDIT_COMMAND_NAME_STATE,
			EditableButton: currButton,
			CurrPath:       callbackData.Path,
			EditableIndex:  callbackData.Index,
			LastMessageID:  int(callbackData.MessageId),
		}
		return keyboard.GetEditCommandNameKeyboard(currButton, callbackData.Index, callbackData.Path)

	case callback_data.ActionType_EDIT_COMMAND_TOPIC:
		user.state = state.StateStruct{
			State:          state.EDIT_COMMAND_TOPIC_STATE,
			EditableButton: currButton,
			CurrPath:       callbackData.Path,
			EditableIndex:  callbackData.Index,
			LastMessageID:  int(callbackData.MessageId),
		}
		return keyboard.GetEditCommandTopicKeyboard(currButton, callbackData.Index, callbackData.Path)

	case callback_data.ActionType_EDIT_COMMAND_VALUE:
		user.state = state.StateStruct{
			State:          state.EDIT_COMMAND_VALUE_STATE,
			EditableButton: currButton,
			CurrPath:       callbackData.Path,
			EditableIndex:  callbackData.Index,
			LastMessageID:  int(callbackData.MessageId),
		}
		return keyboard.GetEditCommandValueKeyboard(currButton, callbackData.Index, callbackData.Path)

	case callback_data.ActionType_SWITCH_SUBSCRIPTION:
		subscriptionButtonIndex := int(callbackData.Index)
		if subscriptionButtonIndex >= len(currButton.GetSubscriptions()) {
			break
		}
		currButton.SetSubscription(subscriptionButtonIndex, int(callbackData.IntValue))
		user.SaveMenuIntoDB()

	case callback_data.ActionType_ADD_SUB_CHART:
		currButton.SetSubscription(len(currButton.GetSubscriptions()), 0)
		user.SaveMenuIntoDB()

	case callback_data.ActionType_DELETE_SUB_CHART:
		if !callbackData.BoolValue {
			return keyboard.GetDeleteSubscriptionKeyboard(currButton, callbackData.Index, user.Subscriptions, callbackData.Path)
		}

		// Yes
		currButton.SetSubscription(int(callbackData.Index), -1)
		user.SaveMenuIntoDB()

	case callback_data.ActionType_BACK_TO_LIST, callback_data.ActionType_BACK_TO_MENU:
		user.state.Reset()
	}

	return keyboard.GetButtonsKeyboard(currButton, callbackData.Path, user.Subscriptions)
}

func (user *User) processChartMenuCallback(callbackData *callback_data.QueryDataType) (string, []byte) {
	var currButton button_interface.ButtonI
	currButton = &user.menu.UserButtons

	for _, buttonIndex := range callbackData.Path {
		if int(buttonIndex) < len(*currButton.GetButtons()) {
			currButton = (*currButton.GetButtons())[buttonIndex]
		}
	}

	var thresholdValue time.Time
	if callbackData.IntValue > 0 {
		thresholdValue = time.Now().Add(-time.Minute * time.Duration(callbackData.IntValue))
	}

	graph := chart.Chart{
		Title: currButton.GetName(),
		TitleStyle: chart.Style{
			Show: true,
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top: 20,
			},
		},
		XAxis: chart.XAxis{
			ValueFormatter: chart.TimeMinuteValueFormatter,
			Style: chart.Style{
				Show: true,
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				Show: true,
			},
		},
		Series: []chart.Series{},
	}

	for index, subscriptionUserIndex := range currButton.GetSubscriptions() {
		if subscriptionUserIndex >= len(user.Subscriptions) {
			continue
		}
		subscription := user.Subscriptions[subscriptionUserIndex]

		user.db.Where("subscription_data.date_time > ?", thresholdValue).
			Order("subscription_data.date_time").Find(&subscription.CollectedData)

		chartTimeSeries := chart.TimeSeries{
			Name: subscription.Topic,
			Style: chart.Style{
				Show:        true,
				StrokeColor: chart.GetDefaultColor(index),
				FillColor:   chart.GetDefaultColor(index).WithAlpha(64),
			},
			XValues: []time.Time{},
			YValues: []float64{},
		}

		for _, subscriptionData := range subscription.CollectedData {
			dataStr := strings.ReplaceAll(string(subscriptionData.Data), ",", ".")
			value, err := strconv.ParseFloat(dataStr, 64)
			if err != nil {
				continue
			}
			chartTimeSeries.XValues = append(chartTimeSeries.XValues, subscriptionData.DateTime)
			chartTimeSeries.YValues = append(chartTimeSeries.YValues, value)
		}

		if len(chartTimeSeries.XValues) > 1 { // there needs to be at least 2 values
			graph.Series = append(graph.Series, chartTimeSeries)
		}
	}

	if len(graph.Series) == 0 {
		return "No data for the selected range", nil
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		err = fmt.Errorf("internal error: could not render chart: %w", err)
		log.Print(err)
		return err.Error(), nil
	}

	return "", buffer.Bytes()
}
