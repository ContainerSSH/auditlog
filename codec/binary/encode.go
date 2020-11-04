package binary

import (
	"compress/gzip"
	"fmt"

	"github.com/containerssh/containerssh-auditlog-go/codec"
	"github.com/containerssh/containerssh-auditlog-go/message"

	"github.com/fxamacker/cbor"
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

func (e *encoder) Encode(messages <-chan message.Message, storage codec.StorageWriter) error {
	var gzipHandle *gzip.Writer
	var encoder *cbor.Encoder
	gzipHandle = gzip.NewWriter(storage)
	encoder = cbor.NewEncoder(gzipHandle, cbor.EncOptions{})
	if err := encoder.StartIndefiniteArray(); err != nil {
		return fmt.Errorf("failed to start infinite array (%v)", err)
	}

	startTime := int64(0)
	var ip = ""
	var username *string
	for {
		msg, ok := <-messages
		if !ok {
			break
		}
		if startTime == 0 {
			startTime = msg.Timestamp
		}
		switch msg.MessageType {
		case message.TypeConnect:
			payload := msg.Payload.(*message.PayloadConnect)
			ip = payload.RemoteAddr
			storage.SetMetadata(startTime/1000000000, ip, username)
		case message.TypeAuthPasswordSuccessful:
			payload := msg.Payload.(*message.PayloadAuthPassword)
			username = &payload.Username
			storage.SetMetadata(startTime/1000000000, ip, username)
		case message.TypeAuthPubKeySuccessful:
			payload := msg.Payload.(*message.PayloadAuthPubKey)
			username = &payload.Username
			storage.SetMetadata(startTime/1000000000, ip, username)
		}
		if err := encoder.Encode(&msg); err != nil {
			return fmt.Errorf("failed to encode audit log message (%v)", err)
		}
		if msg.MessageType == message.TypeDisconnect {
			break
		}
	}
	if err := encoder.EndIndefinite(); err != nil {
		return fmt.Errorf("failed to end audit log infinite array (%v)", err)
	}
	if err := gzipHandle.Flush(); err != nil {
		return fmt.Errorf("failed to flush audit log gzip stream (%v)", err)
	}
	if err := storage.Close(); err != nil {
		return fmt.Errorf("failed to close audit log gzip stream (%v)", err)
	}
	return nil
}
