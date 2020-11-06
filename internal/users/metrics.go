package users

import (
	"github.com/prometheus/client_golang/prometheus"
	"mqtg-bot/internal/common"
)

type Metrics struct {
	common.MetricCollectors
	numOfActiveUsers prometheus.Gauge
	numOfTotalUsers  prometheus.Gauge
}

func InitPrometheusMetrics() Metrics {
	var m Metrics

	m.numOfActiveUsers = m.MetricCollectors.InitMetric("Number of active users")
	m.numOfTotalUsers = m.MetricCollectors.InitMetric("Number of total users")

	return m
}
