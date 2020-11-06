package mqtt

import (
	"github.com/prometheus/client_golang/prometheus"
	"mqtg-bot/internal/common"
)

var metrics struct {
	common.MetricCollectors
	numOfIncMessagesFromMQTT prometheus.Gauge
	numOfOutMessagesToMQTT   prometheus.Gauge
	numOfMqttSubscriptions   prometheus.Gauge
}

func init() {
	metrics.numOfIncMessagesFromMQTT = metrics.MetricCollectors.InitMetric("Number of incoming messages to MQTT")
	metrics.numOfOutMessagesToMQTT = metrics.MetricCollectors.InitMetric("Number of outgoing messages from MQTT")
	metrics.numOfMqttSubscriptions = metrics.MetricCollectors.InitMetric("Number of MQTT subscriptions")
}

func GetPrometheusMetrics() []prometheus.Collector {
	return metrics.Collectors
}
