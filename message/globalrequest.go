package message

type PayloadGlobalRequestUnknown struct {
	RequestType string `json:"requestType" yaml:"requestType"`
}

func (p PayloadGlobalRequestUnknown) Equals(other Payload) bool {
	p2, ok := other.(PayloadGlobalRequestUnknown)
	if !ok {
		return false
	}
	return p.RequestType == p2.RequestType
}
