package message

import "bytes"

type PayloadAuthPassword struct {
	Username string `json:"username" yaml:"username"`
	Password []byte `json:"password" yaml:"password"`
}

func (p * PayloadAuthPassword) Equals(other Payload) bool {
	p2, ok := other.(*PayloadAuthPassword)
	if !ok {
		return false
	}
	return p.Username == p2.Username && bytes.Compare(p.Password, p2.Password) == 0
}

type PayloadAuthPubKey struct {
	Username string `json:"username" yaml:"username"`
	Key      []byte `json:"key" yaml:"key"`
}

func (p * PayloadAuthPubKey) Equals(other Payload) bool {
	p2, ok := other.(*PayloadAuthPubKey)
	if !ok {
		return false
	}
	return p.Username == p2.Username && bytes.Compare(p.Key, p2.Key) == 0
}
