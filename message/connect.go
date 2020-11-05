package message

// PayloadConnect is the payload for TypeConnect messages.
type PayloadConnect struct {
	RemoteAddr string `json:"remoteAddr" yaml:"remoteAddr"`
}

// Equals compares two PayloadConnect datasets.
func (p PayloadConnect) Equals(other Payload) bool {
	p2, ok := other.(PayloadConnect)
	if !ok {
		return false
	}
	return p.RemoteAddr == p2.RemoteAddr
}
