package message

import "bytes"

type ConnectionID []byte
type ChannelID int64

type Message struct {
	// ConnectionID is an opaque ID of the connection
	ConnectionID ConnectionID `json:"connectionId" yaml:"connectionId"`
	// Timestamp is a nanosecond timestamp when the message was created
	Timestamp int64 `json:"timestamp" yaml:"timestamp"`
	// Type of the Payload object
	MessageType MessageType `json:"type" yaml:"type"`
	// Payload is always a pointer to a payload object.
	Payload Payload `json:"payload" yaml:"payload"`
	// ChannelID is a identifier for an SSH channel, if applicable. -1 otherwise.
	ChannelID ChannelID `json:"channelId" yaml:"channelId"`
}

type Payload interface {
	Equals(payload Payload) bool
}

func (m *Message) Equals(other *Message) bool {
	if !bytes.Equal(m.ConnectionID, other.ConnectionID) {
		return false
	}
	if m.Timestamp != other.Timestamp {
		return false
	}
	if m.MessageType != other.MessageType {
		return false
	}
	if m.ChannelID != other.ChannelID {
		return false
	}

	if m.Payload == nil && other.Payload != nil {
		return false
	}
	if m.Payload != nil && other.Payload == nil {
		return false
	}
	if m.Payload != nil {
		return m.Payload.Equals(other.Payload)
	}
	return true
}
