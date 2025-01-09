package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sapcc/bird_exporter/protocol"
)

type LegacyMetricExporter struct {
	ipv4Exporter *GenericProtocolMetricExporter
	ipv6Exporter *GenericProtocolMetricExporter
}

func NewLegacyMetricExporter(prefixIpv4, prefixIpv6 string, labelStrategy LabelStrategy) MetricExporter {
	return &LegacyMetricExporter{
		ipv4Exporter: NewGenericProtocolMetricExporter(prefixIpv4, false, labelStrategy),
		ipv6Exporter: NewGenericProtocolMetricExporter(prefixIpv6, false, labelStrategy),
	}
}

func (e *LegacyMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	e.ipv4Exporter.Describe(ch)
	e.ipv6Exporter.Describe(ch)
}

func (e *LegacyMetricExporter) Export(p *protocol.Protocol, ch chan<- prometheus.Metric, newFormat bool) {
	if p.IPVersion == "4" {
		e.ipv4Exporter.Export(p, ch, false)
	} else {
		e.ipv6Exporter.Export(p, ch, false)
	}
}
