package binary

import (
	"compress/gzip"
	"fmt"
	"io"

	"github.com/containerssh/containerssh-auditlog-go/codec"
	"github.com/containerssh/containerssh-auditlog-go/message"

	"github.com/fxamacker/cbor"
	"github.com/mitchellh/mapstructure"
)

func NewDecoder() codec.Decoder {
	return &decoder{}
}

type decoder struct {
}

func (d *decoder) Decode(reader io.Reader) (<-chan *message.Message, <-chan error, <-chan bool) {
	result := make(chan *message.Message)
	errors := make(chan error)
	done := make(chan bool, 1)

	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		go func() {
			errors <- fmt.Errorf("failed to open gzip stream (%v)", err)
			done <- true
			close(result)
			close(errors)
			close(done)
		}()
		return result, errors, done
	}

	cborReader := cbor.NewDecoder(gzipReader)

	var messages []decodedMessage

	go func() {
		if err = cborReader.Decode(&messages); err != nil {
			errors <- fmt.Errorf("failed to decode messages (%v)", err)
			done <- true
			close(result)
			close(errors)
			close(done)
			return
		}

		for _, v := range messages {
			decodedMessage, err  := decodeMessage(v)
			if err != nil {
				errors <- err
			} else {
				result <- decodedMessage
			}
		}
		done <- true
		close(result)
		close(errors)
		close(done)
	}()
	return result, errors, done
}

type decodedMessage struct {
	// ConnectionID is an opaque ID of the connection
	ConnectionID []byte      `json:"connectionId" yaml:"connectionId"`
	// Timestamp is a nanosecond timestamp when the message was created
	Timestamp    int64       `json:"timestamp" yaml:"timestamp"`
	// Type of the Payload object
	MessageType  message.MessageType `json:"type" yaml:"type"`
	// Payload is always a pointer to a payload object.
	Payload      map[string]interface{} `json:"payload" yaml:"payload"`
	// ChannelID is a identifier for an SSH channel, if applicable. -1 otherwise.
	ChannelID    message.ChannelID   `json:"channelId" yaml:"channelId"`
}

func decodeMessage(v decodedMessage) (*message.Message, error) {
	var payload message.Payload

	switch v.MessageType {

	case message.TypeConnect:
		payload = &message.PayloadConnect{}
	case message.TypeDisconnect:

	case message.TypeAuthPassword:
		payload = &message.PayloadAuthPassword{}
	case message.TypeAuthPasswordSuccessful:
		payload = &message.PayloadAuthPassword{}
	case message.TypeAuthPasswordFailed:
		payload = &message.PayloadAuthPassword{}
	case message.TypeAuthPasswordBackendError:
		payload = &message.PayloadAuthPassword{}

	case message.TypeAuthPubKey:
		payload = &message.PayloadAuthPubKey{}
	case message.TypeAuthPubKeySuccessful:
		payload = &message.PayloadAuthPubKey{}
	case message.TypeAuthPubKeyFailed:
		payload = &message.PayloadAuthPubKey{}
	case message.TypeAuthPubKeyBackendError:
		payload = &message.PayloadAuthPubKey{}

	case message.TypeGlobalRequestUnknown:
		payload = &message.PayloadGlobalRequestUnknown{}
	case message.TypeNewChannel:
		payload = &message.PayloadNewChannel{}
	case message.TypeNewChannelSuccessful:
		payload = &message.PayloadNewChannelSuccessful{}
	case message.TypeNewChannelFailed:
		payload = &message.PayloadNewChannelFailed{}

	case message.TypeChannelRequestUnknownType:
		payload = &message.PayloadChannelRequestUnknownType{}
	case message.TypeChannelRequestDecodeFailed:
		payload = &message.PayloadChannelRequestDecodeFailed{}
	case message.TypeChannelRequestSetEnv:
		payload = &message.PayloadChannelRequestSetEnv{}
	case message.TypeChannelRequestExec:
		payload = &message.PayloadChannelRequestExec{}
	case message.TypeChannelRequestPty:
		payload = &message.PayloadChannelRequestPty{}
	case message.TypeChannelRequestShell:
		payload = &message.PayloadChannelRequestShell{}
	case message.TypeChannelRequestSignal:
		payload = &message.PayloadChannelRequestSignal{}
	case message.TypeChannelRequestSubsystem:
		payload = &message.PayloadChannelRequestSubsystem{}
	case message.TypeChannelRequestWindow:
		payload = &message.PayloadChannelRequestWindow{}
	case message.TypeIO:
		payload = &message.PayloadIO{}
	default:
		return nil, fmt.Errorf("invalid message type: %d", v.MessageType)
	}

	if payload != nil {
		if err := mapstructure.Decode(v.Payload, payload); err != nil {
			return nil, err
		}
	}
	return &message.Message{
		ConnectionID: v.ConnectionID,
		Timestamp:    v.Timestamp,
		MessageType:  v.MessageType,
		Payload:      payload,
		ChannelID:    v.ChannelID,
	}, nil
}