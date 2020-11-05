package message

import "bytes"

// ConnectionID is an opaque, globally unique identifier for a connection made to the SSH server
type ConnectionID []byte

// ChannelID is the ID of an SSH channel
type ChannelID int64

// Message is a basic element of audit logging. It contains the basic records of an interaction.
type Message struct {
	// ConnectionID is an opaque ID of the connection
	ConnectionID ConnectionID `json:"connectionId" yaml:"connectionId"`
	// Timestamp is a nanosecond timestamp when the message was created
	Timestamp int64 `json:"timestamp" yaml:"timestamp"`
	// Type of the Payload object
	MessageType Type `json:"type" yaml:"type"`
	// Payload is always a pointer to a payload object.
	Payload Payload `json:"payload" yaml:"payload"`
	// ChannelID is a identifier for an SSH channel, if applicable. -1 otherwise.
	ChannelID ChannelID `json:"channelId" yaml:"channelId"`
}

// Payload is an interface that makes sure all payloads with Message have a method to compare them.
type Payload interface {
	// Equals compares if the current payload is identical to the provided other payload.
	Equals(payload Payload) bool
}

// Equals is a method to compare two messages with each other
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
