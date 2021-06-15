package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"mqtg-bot/internal/models"
	"net/url"
	"time"
	"os"
)

type Client struct {
	mqtt.Client
	subscriptionCh chan SubscriptionMessage
}

type SubscriptionMessage struct {
	Message      mqtt.Message
	Subscription *models.Subscription
}

func Connect(dbUser *models.DbUser, subscriptionCh chan SubscriptionMessage) (*Client, error) {
	uri, err := url.Parse(dbUser.MqttUrl)
	if err != nil {
		return nil, fmt.Errorf("could not parse MqttUrl: %v", err)
	}

	clientOptions := mqtt.NewClientOptions()
	clientOptions.AddBroker(fmt.Sprintf("%s://%s", uri.Scheme, uri.Host))
	clientOptions.SetUsername(uri.User.Username())
	password, _ := uri.User.Password()
	clientOptions.SetPassword(password)
	clientId := os.Getenv("MQTT_CLIENT_ID")
	if clientId != "" {
		clientOptions.SetClientID(clientId)
	}
	client := mqtt.NewClient(clientOptions)

	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}

	if err := token.Error(); err != nil {
		return nil, err
	}

	return &Client{
		Client:         client,
		subscriptionCh: subscriptionCh,
	}, nil
}

func (c *Client) Subscribe(subscription *models.Subscription) {
	c.Client.Subscribe(subscription.Topic, subscription.Qos, func(client mqtt.Client, msg mqtt.Message) {
		c.subscriptionCh <- SubscriptionMessage{
			Message:      msg,
			Subscription: subscription,
		}
		metrics.numOfIncMessagesFromMQTT.Inc()
	})
	metrics.numOfMqttSubscriptions.Inc()
}

func (c *Client) Unsubscribe(subscription *models.Subscription) {
	c.Client.Unsubscribe(subscription.Topic)
	metrics.numOfMqttSubscriptions.Dec()
}

func (c *Client) Publish(topic string, qos byte, retained bool, payload interface{}) {
	c.Client.Publish(topic, qos, retained, payload)
	metrics.numOfOutMessagesToMQTT.Inc()
}
