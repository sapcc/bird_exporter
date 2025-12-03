package client

import (
	birdsocket "github.com/czerwonk/bird_socket"

	"github.com/sapcc/bird_exporter/parser"
	"github.com/sapcc/bird_exporter/protocol"
)

// BirdClient communicates with the bird socket to retrieve information
type BirdClient struct {
	Options *BirdClientOptions
}

// BirdClientOptions defines options to connect to bird
type BirdClientOptions struct {
	BirdV2       bool
	BirdEnabled  bool
	Bird6Enabled bool
	BirdSocket   string
	Bird6Socket  string
}

// GetProtocols retrieves protocol information and statistics from bird
func (c *BirdClient) GetProtocols() ([]*protocol.Protocol, error) {
	ipVersions := make([]string, 0)
	if c.Options.BirdV2 {
		ipVersions = append(ipVersions, "")
	} else {
		if c.Options.BirdEnabled {
			ipVersions = append(ipVersions, "4")
		}

		if c.Options.Bird6Enabled {
			ipVersions = append(ipVersions, "6")
		}
	}

	return c.protocolsFromBird(ipVersions)
}

// GetOSPFAreas retrieves OSPF specific information from bird
func (c *BirdClient) GetOSPFAreas(p *protocol.Protocol) ([]*protocol.OSPFArea, error) {
	sock := c.socketFor(p.IPVersion)
	b, err := birdsocket.Query(sock, "show ospf "+p.Name)
	if err != nil {
		return nil, err
	}

	return parser.ParseOSPF(b), nil
}

// GetBFDSessions retrieves BFD specific information from bird
func (c *BirdClient) GetBFDSessions(p *protocol.Protocol) ([]*protocol.BFDSession, error) {
	sock := c.socketFor(p.IPVersion)
	b, err := birdsocket.Query(sock, "show bfd sessions "+p.Name)
	if err != nil {
		return nil, err
	}

	return parser.ParseBFDSessions(p.Name, b), nil
}

func (c *BirdClient) protocolsFromBird(ipVersions []string) ([]*protocol.Protocol, error) {
	protocols := make([]*protocol.Protocol, 0)

	for _, ipVersion := range ipVersions {
		sock := c.socketFor(ipVersion)
		s, err := c.protocolsFromSocket(sock, ipVersion)
		if err != nil {
			return nil, err
		}

		protocols = append(protocols, s...)
	}

	return protocols, nil
}

func (c *BirdClient) protocolsFromSocket(socketPath, ipVersion string) ([]*protocol.Protocol, error) {
	b, err := birdsocket.Query(socketPath, "show protocols all")
	if err != nil {
		return nil, err
	}

	return parser.ParseProtocols(b, ipVersion), nil
}

func (c *BirdClient) socketFor(ipVersion string) string {
	if !c.Options.BirdV2 && ipVersion == "6" {
		return c.Options.Bird6Socket
	}

	return c.Options.BirdSocket
}
