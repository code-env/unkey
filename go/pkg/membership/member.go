package membership

import (
	"encoding/json"
	"net"
)

type Member struct {
	// Global unique identifier for the node
	NodeID string `json:"nodeId"`
	Addr   net.IP `json:"addr"`
}

func (m Member) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Member) Unmarshal(b []byte) error {
	return json.Unmarshal(b, m)
}
