package keyboard

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"mqtg-bot/internal/models"
	"mqtg-bot/internal/users/keyboard/callback_data"
	"mqtg-bot/internal/users/keyboard/keyboard_names"
)

func GetSubscriptionsKeyboard(subscriptions []*models.Subscription) (string, *tgbotapi.InlineKeyboardMarkup) {
	var text = fmt.Sprintf("You have %v subscriptions:", len(subscriptions))

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup()
	for index, subscription := range subscriptions {
		inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				subscription.Topic,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_SUBSCRIPTIONS,
					Path:     []int32{int32(index)},
				}.GetBase64ProtoString()),
		))
	}

	inlineKeyboard.InlineKeyboard = append(inlineKeyboard.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			keyboard_names.ADD_SUBSCRIPTION,
			callback_data.QueryDataType{
				Keyboard: callback_data.KeyboardType_SUBSCRIPTIONS,
				Action:   callback_data.ActionType_ADD_SUBSCRIPTION,
			}.GetBase64ProtoString()),
	))

	return text, &inlineKeyboard
}

func GetAddSubscriptionKeyboard(subType models.SubscriptionType, qos byte, subDataType models.SubscriptionDataType) (string, *tgbotapi.InlineKeyboardMarkup) {
	inlineText := "Select subscription type, QoS level and subscription data type"
	inlineText += "\nThen send me the <code>topic</code> for subscribing"
	inlineText += "\n  <code>+</code> (plus symbol) represents the single-level wildcard in the topic"
	inlineText += "\n  <code>#</code> (hash symbol) represents the multi-level wildcard in the topic"
	inlineText += "\nFor example: <code>/sensors/+/temp</code>"

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(keyboard_names.SUBSCRIPTION_TYPE, subType),
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_SUBSCRIPTIONS,
					Action:   callback_data.ActionType_SWITCH_ADDED_SUBSCRIPTION_TYPE,
					IntValue: int32(subType.GetNext()),
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(keyboard_names.QOS, qos),
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_SUBSCRIPTIONS,
					Action:   callback_data.ActionType_SWITCH_ADDED_SUBSCRIPTION_QOS,
					IntValue: int32(qos+1) % 3,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(keyboard_names.DATA_TYPE, subDataType),
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_SUBSCRIPTIONS,
					Action:   callback_data.ActionType_SWITCH_ADDED_SUBSCRIPTION_DATA_TYPE,
					IntValue: int32(subDataType.GetNext()),
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.BACK_TO_LIST,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_SUBSCRIPTIONS,
					Action:   callback_data.ActionType_BACK_TO_LIST,
				}.GetBase64ProtoString()),
		),
	)

	return inlineText, &inlineKeyboard
}

func GetSubscriptionTopicEditKeyboard(subscription *models.Subscription, subscriptionIndex int32) (string, *tgbotapi.InlineKeyboardMarkup) {
	inlineText := fmt.Sprintf("Current topic is: <code>%v</code>", subscription.Topic)
	inlineText += "\nSend me the new topic if you want to edit it"
	inlineText += "\n  <code>+</code> (plus symbol) represents the single-level wildcard in the topic"
	inlineText += "\n  <code>#</code> (hash symbol) represents the multi-level wildcard in the topic"
	inlineText += "\nFor example: <code>/sensors/+/temp</code>"

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.BACK_TO_MENU,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_SUBSCRIPTIONS,
					Path:     []int32{subscriptionIndex},
					Action:   callback_data.ActionType_BACK_TO_MENU,
				}.GetBase64ProtoString()),
		),
	)

	return inlineText, &inlineKeyboard
}

func GetSubscriptionDeleteKeyboard(subscription *models.Subscription, subscriptionIndex int32) (string, *tgbotapi.InlineKeyboardMarkup) {
	inlineText := fmt.Sprintf("Are you sure you want to delete the subscription <code>%v</code>?", subscription.Topic)
	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.YES,
				callback_data.QueryDataType{
					Keyboard:  callback_data.KeyboardType_SUBSCRIPTIONS,
					Path:      []int32{subscriptionIndex},
					Action:    callback_data.ActionType_DELETE,
					BoolValue: true,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.BACK,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_SUBSCRIPTIONS,
					Path:     []int32{subscriptionIndex},
					Action:   callback_data.ActionType_BACK_TO_MENU,
				}.GetBase64ProtoString()),
		),
	)

	return inlineText, &inlineKeyboard
}

func GetSubscriptionBeforeAfterValueTextEditKeyboard(action callback_data.ActionType, subscription *models.Subscription, subscriptionIndex int32) (string, *tgbotapi.InlineKeyboardMarkup) {
	var inlineText string
	if action == callback_data.ActionType_BEFORE_VALUE_TEXT {
		if subscription.BeforeValueText == "" {
			inlineText = fmt.Sprintf("Current <code>before-value</code> message text is empty")
		} else {
			inlineText = fmt.Sprintf("Current <code>before-value</code> message text is: <code>%v</code>", subscription.BeforeValueText)
		}
	} else {
		if subscription.AfterValueText == "" {
			inlineText = fmt.Sprintf("Current <code>after-value</code> message text is empty")
		} else {
			inlineText = fmt.Sprintf("Current <code>after-value</code> message text is: <code>%v</code>", subscription.AfterValueText)
		}
	}

	inlineText += "\nSend me the new text if you want to edit it"
	inlineText += "\nYou can use:"
	inlineText += "\n  <code>%e</code> for empty string"
	inlineText += "\n  <code>%s</code> for printing the subscription topic as is (with + or #)"
	inlineText += "\n  <code>%t</code> for printing the full-level received topic"

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.BACK_TO_MENU,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_SUBSCRIPTIONS,
					Path:     []int32{subscriptionIndex},
					Action:   callback_data.ActionType_BACK_TO_MENU,
				}.GetBase64ProtoString()),
		),
	)

	return inlineText, &inlineKeyboard
}

func GetSubscriptionEditKeyboard(subscription *models.Subscription, subscriptionIndex int32) (string, *tgbotapi.InlineKeyboardMarkup) {
	inlineText := fmt.Sprintf("Edit the subscription <code>%v</code>:", subscription.Topic)

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(keyboard_names.SUBSCRIPTION_TYPE, subscription.SubscriptionType),
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_SUBSCRIPTIONS,
					Path:     []int32{subscriptionIndex},
					Action:   callback_data.ActionType_SWITCH_SUBSCRIPTION_TYPE,
					IntValue: int32(subscription.SubscriptionType.GetNext()),
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(keyboard_names.QOS, subscription.Qos),
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_SUBSCRIPTIONS,
					Path:     []int32{subscriptionIndex},
					Action:   callback_data.ActionType_SWITCH_QOS,
					IntValue: int32(subscription.Qos+1) % 3,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf(keyboard_names.DATA_TYPE, subscription.DataType),
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_SUBSCRIPTIONS,
					Path:     []int32{subscriptionIndex},
					Action:   callback_data.ActionType_SWITCH_SUB_DATA_TYPE,
					IntValue: int32(subscription.DataType.GetNext()),
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.BEFORE_VALUE_MESSAGE_TEXT,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_SUBSCRIPTIONS,
					Path:     []int32{subscriptionIndex},
					Action:   callback_data.ActionType_BEFORE_VALUE_TEXT,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.AFTER_VALUE_MESSAGE_TEXT,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_SUBSCRIPTIONS,
					Path:     []int32{subscriptionIndex},
					Action:   callback_data.ActionType_AFTER_VALUE_TEXT,
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.EDIT_TOPIC,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_SUBSCRIPTIONS,
					Path:     []int32{subscriptionIndex},
					Action:   callback_data.ActionType_EDIT,
				}.GetBase64ProtoString()),
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.DELETE_SUBSCRIPTION,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_SUBSCRIPTIONS,
					Path:     []int32{subscriptionIndex},
					Action:   callback_data.ActionType_DELETE,
				}.GetBase64ProtoString()),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				keyboard_names.BACK_TO_LIST,
				callback_data.QueryDataType{
					Keyboard: callback_data.KeyboardType_SUBSCRIPTIONS,
					Action:   callback_data.ActionType_BACK_TO_LIST,
				}.GetBase64ProtoString()),
		),
	)

	return inlineText, &inlineKeyboard
}
