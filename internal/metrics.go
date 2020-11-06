package internal

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"mqtg-bot/internal/common"
	"net/http"
	_ "net/http/pprof"
	"os"
)

const DEFAULT_PORT = "80"

type Metrics struct {
	common.MetricCollectors

	numOfIncMessagesFromTelegram prometheus.Gauge
	numOfOutMessagesToTelegram   prometheus.Gauge
}

func InitPrometheusMetrics() Metrics {
	var m Metrics

	m.numOfIncMessagesFromTelegram = m.MetricCollectors.InitMetric("Number of incoming messages from Telegram")
	m.numOfOutMessagesToTelegram = m.MetricCollectors.InitMetric("Number of outgoing messages to Telegram")

	return m
}

func (bot *TelegramBot) StartPprofAndMetricsListener() {
	http.Handle("/metrics", promhttp.Handler())

	port := os.Getenv("PORT")
	if port == "" {
		port = DEFAULT_PORT
	}

	http.ListenAndServe(":"+port, nil)
}
