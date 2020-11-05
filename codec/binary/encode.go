package binary

import (
	"compress/gzip"
	"fmt"

	"github.com/containerssh/auditlog/storage"

	"github.com/containerssh/auditlog/codec"
	"github.com/containerssh/auditlog/message"

	"github.com/fxamacker/cbor"
)

// NewEncoder creates an encoder that encodes messages in CBOR+GZIP format as documented
//            on https://containerssh.github.io/advanced/audit/format/
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

func (e *encoder) Encode(messages <-chan message.Message, storage storage.Writer) error {
	var gzipHandle *gzip.Writer
	var encoder *cbor.Encoder
	gzipHandle = gzip.NewWriter(storage)
	encoder = cbor.NewEncoder(gzipHandle, cbor.EncOptions{})
	if err := encoder.StartIndefiniteArray(); err != nil {
		return fmt.Errorf("failed to start infinite array (%w)", err)
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
			payload := msg.Payload.(message.PayloadConnect)
			ip = payload.RemoteAddr
			storage.SetMetadata(startTime/1000000000, ip, username)
		case message.TypeAuthPasswordSuccessful:
			payload := msg.Payload.(message.PayloadAuthPassword)
			username = &payload.Username
			storage.SetMetadata(startTime/1000000000, ip, username)
		case message.TypeAuthPubKeySuccessful:
			payload := msg.Payload.(message.PayloadAuthPubKey)
			username = &payload.Username
			storage.SetMetadata(startTime/1000000000, ip, username)
		}
		if err := encoder.Encode(&msg); err != nil {
			return fmt.Errorf("failed to encode audit log message (%w)", err)
		}
		if msg.MessageType == message.TypeDisconnect {
			break
		}
	}
	if err := encoder.EndIndefinite(); err != nil {
		return fmt.Errorf("failed to end audit log infinite array (%w)", err)
	}
	if err := gzipHandle.Flush(); err != nil {
		return fmt.Errorf("failed to flush audit log gzip stream (%w)", err)
	}
	if err := storage.Close(); err != nil {
		return fmt.Errorf("failed to close audit log gzip stream (%w)", err)
	}
	return nil
}
