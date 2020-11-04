package message

import "bytes"

type PayloadAuthPassword struct {
	Username string `json:"username" yaml:"username"`
	Password []byte `json:"password" yaml:"password"`
}

func (p *PayloadAuthPassword) Equals(other Payload) bool {
	p2, ok := other.(*PayloadAuthPassword)
	if !ok {
		return false
	}
	return p.Username == p2.Username && bytes.Compare(p.Password, p2.Password) == 0
}

type PayloadAuthPasswordBackendError struct {
	Username string `json:"username" yaml:"username"`
	Password []byte `json:"password" yaml:"password"`
	Reason   string `json:"reason" yaml:"reason"`
}

func (p *PayloadAuthPasswordBackendError) Equals(other Payload) bool {
	p2, ok := other.(*PayloadAuthPasswordBackendError)
	if !ok {
		return false
	}
	return p.Username == p2.Username && bytes.Compare(p.Password, p2.Password) == 0 && p.Reason == p2.Reason
}

type PayloadAuthPubKey struct {
	Username string `json:"username" yaml:"username"`
	Key      []byte `json:"key" yaml:"key"`
}

func (p *PayloadAuthPubKey) Equals(other Payload) bool {
	p2, ok := other.(*PayloadAuthPubKey)
	if !ok {
		return false
	}
	return p.Username == p2.Username && bytes.Compare(p.Key, p2.Key) == 0
}

type PayloadAuthPubKeyBackendError struct {
	Username string `json:"username" yaml:"username"`
	Key      []byte `json:"key" yaml:"key"`
	Reason   string `json:"reason" yaml:"reason"`
}

func (p *PayloadAuthPubKeyBackendError) Equals(other Payload) bool {
	p2, ok := other.(*PayloadAuthPubKeyBackendError)
	if !ok {
		return false
	}
	return p.Username == p2.Username && bytes.Compare(p.Key, p2.Key) == 0 && p.Reason == p2.Reason
}
