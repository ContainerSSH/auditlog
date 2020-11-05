package message

import "bytes"

// PayloadAuthPassword is a payload for a message that indicates an authentication attempt, successful, or failed
//                     authentication.
type PayloadAuthPassword struct {
	Username string `json:"username" yaml:"username"`
	Password []byte `json:"password" yaml:"password"`
}

// Equals compares two PayloadAuthPassword payloads.
func (p PayloadAuthPassword) Equals(other Payload) bool {
	p2, ok := other.(PayloadAuthPassword)
	if !ok {
		return false
	}
	return p.Username == p2.Username && bytes.Equal(p.Password, p2.Password)
}

// PayloadAuthPasswordBackendError is a payload for a message that indicates a backend failure during authentication.
type PayloadAuthPasswordBackendError struct {
	Username string `json:"username" yaml:"username"`
	Password []byte `json:"password" yaml:"password"`
	Reason   string `json:"reason" yaml:"reason"`
}

// Equals compares two PayloadAuthPasswordBackendError payloads.
func (p PayloadAuthPasswordBackendError) Equals(other Payload) bool {
	p2, ok := other.(PayloadAuthPasswordBackendError)
	if !ok {
		return false
	}
	return p.Username == p2.Username && bytes.Equal(p.Password, p2.Password) && p.Reason == p2.Reason
}

// PayloadAuthPubKey is a payload for a public key based authentication
type PayloadAuthPubKey struct {
	Username string `json:"username" yaml:"username"`
	Key      []byte `json:"key" yaml:"key"`
}

// Equals compares two PayloadAuthPubKey payloads
func (p PayloadAuthPubKey) Equals(other Payload) bool {
	p2, ok := other.(PayloadAuthPubKey)
	if !ok {
		return false
	}
	return p.Username == p2.Username && bytes.Equal(p.Key, p2.Key)
}

// PayloadAuthPubKeyBackendError is a payload for a message indicating that there was a backend error while
//                               authenticating with public key.
type PayloadAuthPubKeyBackendError struct {
	Username string `json:"username" yaml:"username"`
	Key      []byte `json:"key" yaml:"key"`
	Reason   string `json:"reason" yaml:"reason"`
}

// Equals compares two PayloadAuthPubKeyBackendError payloads
func (p PayloadAuthPubKeyBackendError) Equals(other Payload) bool {
	p2, ok := other.(PayloadAuthPubKeyBackendError)
	if !ok {
		return false
	}
	return p.Username == p2.Username && bytes.Equal(p.Key, p2.Key) && p.Reason == p2.Reason
}
