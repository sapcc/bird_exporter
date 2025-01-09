package metrics

import "github.com/sapcc/bird_exporter/protocol"

// LabelStrategy abstracts the label generation for protocol metrics
type LabelStrategy interface {
	// LabelNames is the list of label names
	LabelNames(p *protocol.Protocol) []string

	// Label values is the list of values for the labels specified in `LabelNames()`
	LabelValues(p *protocol.Protocol) []string
}
