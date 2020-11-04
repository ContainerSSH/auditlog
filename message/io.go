package message

import "bytes"

type Stream uint

const (
	StreamStdin  Stream = 0
	StreamStdout Stream = 1
	StreamStderr Stream = 2
)

type PayloadIO struct {
	Stream Stream `json:"stream" yaml:"stream"`
	Data   []byte `json:"data" yaml:"data"`
}

func (p *PayloadIO) Equals(other Payload) bool {
	p2, ok := other.(*PayloadIO)
	if !ok {
		return false
	}
	return p.Stream == p2.Stream && bytes.Compare(p.Data, p2.Data) == 0
}
