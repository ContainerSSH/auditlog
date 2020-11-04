package binary

import (
	"github.com/containerssh/containerssh-auditlog-go/codec"
	"github.com/containerssh/containerssh-auditlog-go/message"
)

func NewEncoder() codec.Encoder {
	return &encoder{}
}

type encoder struct {
}

func (e *encoder) GetMimeType() string {
	return "application/octet-stream"
}

func (e *encoder) GetFileExtension() string {
	return ""
}

func (e *encoder) Encode(messages <-chan message.Message, _ codec.StorageWriter) error {
	for {
		msg, ok := <-messages
		if !ok {
			break
		}
		if msg.MessageType == message.TypeDisconnect {
			break
		}
	}
	return nil
}
