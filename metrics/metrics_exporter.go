package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sapcc/bird_exporter/protocol"
)

type MetricExporter interface {
	Describe(ch chan<- *prometheus.Desc)
	Export(p *protocol.Protocol, ch chan<- prometheus.Metric, newFormat bool)
}
