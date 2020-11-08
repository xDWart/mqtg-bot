package users

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"mqtg-bot/internal/common"
	"mqtg-bot/internal/models"
	"mqtg-bot/internal/users/keyboard"
	"mqtg-bot/internal/users/menu"
	"mqtg-bot/internal/users/menu/button_interface"
	"mqtg-bot/internal/users/menu/button_names"
	"mqtg-bot/internal/users/menu/button_types"
	"mqtg-bot/internal/users/state"
	"strings"
)

func (user *User) ProcessMessage(messageData []byte, isItPhoto bool) *common.BotMessage {
	if len(messageData) == 0 {
		return nil
	}

	if isItPhoto && user.state.State != state.PUBLISH_VALUE_STATE {
		return nil // ignore it
	}

	// the users buttons names must not match with the system buttons names
	switch string(messageData) {
	case button_names.PUBLISH,
		button_names.SETTINGS,
		button_names.BACK,
		button_names.SUBSCRIPTIONS,
		button_names.EDIT_BUTTONS,
		button_names.CONFIGURE_CONNECTION,
		button_names.DISCONNECT,
		button_names.START:
		user.state.Reset()
	}

	switch user.state.State {
	case state.PUBLISH_TOPIC_STATE:
		user.state = state.StateStruct{
			State:           state.PUBLISH_VALUE_STATE,
			PublishingTopic: string(messageData),
		}

		return &common.BotMessage{
			MainText: "Ok, now send me the payload data for publish to this topic",
		}

	case state.PUBLISH_VALUE_STATE:
		defer user.state.Reset()

		user.publish(messageData)
		user.menu.Back()
		return &common.BotMessage{
			MainText: fmt.Sprintf("Successfully published to the topic <code>%v</code>", user.state.PublishingTopic),
			MainMenu: user.getMainMenu(),
		}

	case state.ADD_SUBSCRIPTION_STATE:
		defer user.state.Reset()

		newSubscription := &models.Subscription{
			UserMutex:        &user.Mutex,
			DbUserID:         user.ID,
			ChatID:           user.ChatID,
			BeforeValueText:  "%t",
			Topic:            string(messageData),
			Qos:              user.state.Qos,
			SubscriptionType: user.state.SubscriptionType,
			DataType:         user.state.SubscriptionDataType,
		}

		subscriptionIndex := user.subscribe(newSubscription)

		inlineText, inlineKeyboard := keyboard.GetSubscriptionEditKeyboard(newSubscription, subscriptionIndex)

		return &common.BotMessage{
			InlineText:     inlineText,
			InlineKeyboard: inlineKeyboard,
		}

	case state.EDIT_BEFORE_VALUE_MESSAGE_TEXT_STATE, state.EDIT_AFTER_VALUE_MESSAGE_TEXT_STATE, state.EDIT_SUBSCRIPTION_TOPIC_STATE:
		defer user.state.Reset()

		subscriptionIndex := user.state.EditableIndex
		if int(subscriptionIndex) >= len(user.Subscriptions) {
			return &common.BotMessage{
				MainText: fmt.Sprintf("Cannot get subscription by index %v", subscriptionIndex),
			}
		}

		message := string(messageData)
		subscription := user.Subscriptions[subscriptionIndex]

		switch user.state.State {
		case state.EDIT_BEFORE_VALUE_MESSAGE_TEXT_STATE:
			subscription.BeforeValueText = strings.ReplaceAll(message, "%e", "")
		case state.EDIT_AFTER_VALUE_MESSAGE_TEXT_STATE:
			subscription.AfterValueText = strings.ReplaceAll(message, "%e", "")
		case state.EDIT_SUBSCRIPTION_TOPIC_STATE:
			user.mqtt.Unsubscribe(subscription)
			subscription.Topic = message
			user.mqtt.Subscribe(subscription)
		}

		user.db.Save(subscription)

		newText, newInlineKeyboard := keyboard.GetSubscriptionEditKeyboard(subscription, subscriptionIndex)

		return &common.BotMessage{
			MessageID:      user.state.LastMessageID,
			InlineText:     newText,
			InlineKeyboard: newInlineKeyboard,
		}

	case state.ADD_BUTTON_STATE,
		state.EDIT_BUTTON_NAME_STATE,
		state.ADD_NEW_COMMAND,
		state.EDIT_COMMAND_NAME_STATE,
		state.EDIT_COMMAND_TOPIC_STATE,
		state.EDIT_COMMAND_VALUE_STATE:

		defer user.state.Reset()
		var newInlineText string
		var newInlineKeyboard *tgbotapi.InlineKeyboardMarkup
		var messageStr = string(messageData)
		var currButton = user.state.EditableButton

		switch user.state.State {
		case state.ADD_BUTTON_STATE:
			newButton, err := button_interface.GetNewButtonWithName(user.state.ButtonType, messageStr)
			if err != nil {
				return &common.BotMessage{
					MainText: err.Error(),
				}
			}
			currButton.AddButton(newButton)
			buttonPath := append(user.state.CurrPath, int32(len(*currButton.GetButtons())-1))
			newInlineText, newInlineKeyboard = keyboard.GetButtonsKeyboard(newButton, buttonPath, user.Subscriptions)
		case state.EDIT_BUTTON_NAME_STATE:
			currButton.SetMainName(messageStr)
			newInlineText, newInlineKeyboard = keyboard.GetButtonsKeyboard(currButton, user.state.CurrPath, user.Subscriptions)
		case state.ADD_NEW_COMMAND:
			currButton.AddNewCommand(&button_interface.CommandType{
				Name:     string(messageData),
				Qos:      user.state.Qos,
				Retained: user.state.Retained,
			})
			newInlineText, newInlineKeyboard = keyboard.GetButtonsKeyboard(currButton, user.state.CurrPath, user.Subscriptions)
		case state.EDIT_COMMAND_NAME_STATE:
			currButton.SetNameForCommand(int(user.state.EditableIndex), string(messageData))
			newInlineText, newInlineKeyboard = keyboard.GetCommandEditKeyboard(currButton, user.state.EditableIndex, user.state.CurrPath)
		case state.EDIT_COMMAND_TOPIC_STATE:
			currButton.SetTopicForCommand(int(user.state.EditableIndex), string(messageData))
			newInlineText, newInlineKeyboard = keyboard.GetCommandEditKeyboard(currButton, user.state.EditableIndex, user.state.CurrPath)
		case state.EDIT_COMMAND_VALUE_STATE:
			currButton.SetValueForCommand(int(user.state.EditableIndex), string(messageData))
			newInlineText, newInlineKeyboard = keyboard.GetCommandEditKeyboard(currButton, user.state.EditableIndex, user.state.CurrPath)
		}

		user.SaveMenuIntoDB()

		return &common.BotMessage{
			MessageID:      user.state.LastMessageID,
			InlineText:     newInlineText,
			InlineKeyboard: newInlineKeyboard,
		}

	default: // it could be any other button
		message := string(messageData)

		if !user.isMqttConnected() {
			return user.processConnectionString(message)
		}

		// pre action
		switch message {
		case button_names.DISCONNECT:
			user.disconnectMQTT()
			return &common.BotMessage{
				MainText: "Disconnected from MQTT broker",
				MainMenu: &menu.ConfigureConnectionMenu,
			}
		case button_names.EDIT_BUTTONS:
			inlineText, inlineKeyboard := keyboard.GetButtonsKeyboard(&user.menu.UserButtons, []int32{}, user.Subscriptions)
			return &common.BotMessage{
				InlineText:     inlineText,
				InlineKeyboard: inlineKeyboard,
			}
		case button_names.SUBSCRIPTIONS:
			inlineText, inlineKeyboard := keyboard.GetSubscriptionsKeyboard(user.Subscriptions)
			return &common.BotMessage{
				InlineText:     inlineText,
				InlineKeyboard: inlineKeyboard,
			}
		case button_names.PUBLISH:
			inlineText, inlineKeyboard := keyboard.GetPublishKeyboard(0, false)
			user.state = state.StateStruct{
				State: state.PUBLISH_TOPIC_STATE,
			}

			return &common.BotMessage{
				InlineText:     inlineText,
				InlineKeyboard: inlineKeyboard,
			}
		}

		// finding button
		button := user.menu.SetPressedButtonLikeCurrentPath(message, user.menu.CurrPath)
		if button == nil {
			return &common.BotMessage{
				MainMenu: user.menu.GetCurrPathMainMenu(),
				MainText: user.menu.CurrPath.GetName(),
			}
		}

		// action
		switch button.GetType() {
		case button_types.FOLDER:
			return &common.BotMessage{
				MainMenu: user.menu.GetCurrPathMainMenu(),
				MainText: user.menu.CurrPath.GetName(),
			}
		case button_types.SINGLE_VALUE, button_types.TOGGLE:
			mainText := button.GetName()

			command := button.GetCurrentCommand()
			if command != nil {
				if len(command.Topic) > 0 && len(command.Value) > 0 {
					user.mqtt.Publish(command.Topic, command.Qos, command.Retained, command.Value)
				} else {
					mainText += " - nothing to send, you need to configure command for this button"
				}
			}

			button.SwitchState()

			return &common.BotMessage{
				MainMenu: user.menu.GetCurrPathMainMenu(),
				MainText: mainText,
			}

		case button_types.MULTI_VALUE:
			inlineText, inlineKeyboard := keyboard.GetMultiValueCommandsKeyboard(button)

			return &common.BotMessage{
				InlineText:     inlineText,
				InlineKeyboard: inlineKeyboard,
			}

		case button_types.PRINT_LAST_SUB_VALUE:
			var userAnswer = common.BotMessage{
				MainMenu: user.menu.GetCurrPathMainMenu(),
			}

			var subscriptionID = button.GetSubscriptions()[0]
			if subscriptionID == -1 {
				userAnswer.MainText = "Subscription not selected, configure it in button settings"
				return &userAnswer
			}

			if subscriptionID < len(user.Subscriptions) {
				subscription := user.Subscriptions[subscriptionID]
				userAnswer.MainText = subscription.LastValueFormattedMessage
				if subscription.DataType == models.IMAGE_DATA_TYPE {
					userAnswer.Photo = subscription.LastValuePayload
				}
			} else {
				userAnswer.MainText = "Subscription not found, configure it in button settings"
			}

			return &userAnswer

		case button_types.DRAW_CHART:
			var userAnswer common.BotMessage

			currentPath := make([]int32, 0)
			currButton := button
			parent := currButton.GetParent()
			for parent != nil {
				for index, iButton := range *parent.GetButtons() {
					if iButton == currButton {
						currentPath = append([]int32{int32(index)}, currentPath...)
						continue
					}
				}
				currButton = parent
				parent = currButton.GetParent()
			}

			userAnswer.InlineText, userAnswer.InlineKeyboard = keyboard.GetShowChartKeyboard(button, currentPath)

			return &userAnswer
		}
	}
	return nil
}
