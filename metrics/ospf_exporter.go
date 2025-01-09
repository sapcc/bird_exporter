package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/sapcc/bird_exporter/client"
	"github.com/sapcc/bird_exporter/protocol"
)

type ospfDesc struct {
	runningDesc               *prometheus.Desc
	interfaceCountDesc        *prometheus.Desc
	neighborCountDesc         *prometheus.Desc
	neighborAdjacentCountDesc *prometheus.Desc
}

type ospfMetricExporter struct {
	descriptions map[string]*ospfDesc
	client       client.Client
}

// NewOSPFExporter creates a new MetricExporter for OSPF metrics
func NewOSPFExporter(prefix string, client client.Client) MetricExporter {
	d := make(map[string]*ospfDesc)
	d["4"] = getDesc(prefix + "ospf")
	d["6"] = getDesc(prefix + "ospfv3")

	return &ospfMetricExporter{descriptions: d, client: client}
}

func getDesc(prefix string) *ospfDesc {
	labels := []string{"name"}

	d := &ospfDesc{}
	d.runningDesc = prometheus.NewDesc(prefix+"_running", "State of OSPF: 0 = Alone, 1 = Running (Neighbor-Adjacencies established)", labels, nil)

	labels = append(labels, "area")
	d.interfaceCountDesc = prometheus.NewDesc(prefix+"_interface_count", "Number of interfaces in the area", labels, nil)
	d.neighborCountDesc = prometheus.NewDesc(prefix+"_neighbor_count", "Number of neighbors in the area", labels, nil)
	d.neighborAdjacentCountDesc = prometheus.NewDesc(prefix+"_neighbor_adjacent_count", "Number of adjacent neighbors in the area", labels, nil)

	return d
}

func (m *ospfMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	m.describe("4", ch)
	m.describe("6", ch)
}

func (m *ospfMetricExporter) describe(ipVersion string, ch chan<- *prometheus.Desc) {
	d := m.descriptions[ipVersion]
	ch <- d.runningDesc
	ch <- d.interfaceCountDesc
	ch <- d.neighborCountDesc
	ch <- d.neighborAdjacentCountDesc
}

func (m *ospfMetricExporter) Export(p *protocol.Protocol, ch chan<- prometheus.Metric, newFormat bool) {
	d := m.descriptions[p.IPVersion]

	var running float64
	if p.State == "Running" {
		running = 1
	}

	ch <- prometheus.MustNewConstMetric(d.runningDesc, prometheus.GaugeValue, running, p.Name)

	areas, err := m.client.GetOSPFAreas(p)
	if err != nil {
		log.Errorln(err)
		return
	}

	for _, area := range areas {
		l := []string{p.Name, area.Name}
		ch <- prometheus.MustNewConstMetric(d.interfaceCountDesc, prometheus.GaugeValue, float64(area.InterfaceCount), l...)
		ch <- prometheus.MustNewConstMetric(d.neighborCountDesc, prometheus.GaugeValue, float64(area.NeighborCount), l...)
		ch <- prometheus.MustNewConstMetric(d.neighborAdjacentCountDesc, prometheus.GaugeValue, float64(area.NeighborAdjacentCount), l...)
	}
}
