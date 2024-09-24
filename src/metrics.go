package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var RequestCounterMetric = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "app_request_count",
		Help: "The total number of request",
	},
	[]string{"node"},
)

var ResponseStatusCodeCounterMetric = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "app_status_code_counter",
		Help: "The status code counter",
	},
	[]string{"node", "path", "status_code"},
)

var ConfigVersionGaugeMetric = promauto.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "app_config_version_gauge",
		Help: "The config version gauge",
	},
	[]string{"config_version"},
)
