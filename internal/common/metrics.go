package common

import (
	"github.com/prometheus/client_golang/prometheus"
	"strings"
)

type MetricCollectors struct {
	Collectors []prometheus.Collector
}

func (metrics *MetricCollectors) InitMetric(help string) prometheus.Gauge {
	promGauge := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: strings.Replace(strings.ToLower(help), " ", "_", -1),
			Help: help,
		},
	)
	metrics.Collectors = append(metrics.Collectors, promGauge)
	return promGauge
}

func (metrics *MetricCollectors) GetPrometheusMetrics() []prometheus.Collector {
	return metrics.Collectors
}
