package binary

import (
	"compress/gzip"
	"fmt"
	"io"

	"github.com/containerssh/auditlog/codec"
	"github.com/containerssh/auditlog/message"

	"github.com/fxamacker/cbor"
	"github.com/mitchellh/mapstructure"
)

// NewDecoder Creates a decoder for the CBOR+GZIP audit log format.
func NewDecoder() codec.Decoder {
	return &decoder{}
}

type decoder struct {
}

func (d *decoder) Decode(reader io.Reader) (<-chan message.Message, <-chan error, <-chan bool) {
	result := make(chan message.Message)
	errors := make(chan error)
	done := make(chan bool, 1)

	gzipReader, err := gzip.NewReader(reader)
	if err != nil {
		go func() {
			errors <- fmt.Errorf("failed to open gzip stream (%w)", err)
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
			errors <- fmt.Errorf("failed to decode messages (%w)", err)
			done <- true
			close(result)
			close(errors)
			close(done)
			return
		}

		for _, v := range messages {
			decodedMessage, err := decodeMessage(v)
			if err != nil {
				errors <- err
			} else {
				result <- *decodedMessage
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
	ConnectionID []byte `json:"connectionId" yaml:"connectionId"`
	// Timestamp is a nanosecond timestamp when the message was created
	Timestamp int64 `json:"timestamp" yaml:"timestamp"`
	// Type of the Payload object
	MessageType message.Type `json:"type" yaml:"type"`
	// Payload is always a pointer to a payload object.
	Payload map[string]interface{} `json:"payload" yaml:"payload"`
	// ChannelID is a identifier for an SSH channel, if applicable. -1 otherwise.
	ChannelID message.ChannelID `json:"channelId" yaml:"channelId"`
}

var messageTypeMap = map[message.Type]message.Payload{
	message.TypeConnect:    message.PayloadConnect{},
	message.TypeDisconnect: nil,

	message.TypeAuthPassword:             message.PayloadAuthPassword{},
	message.TypeAuthPasswordSuccessful:   message.PayloadAuthPassword{},
	message.TypeAuthPasswordFailed:       message.PayloadAuthPassword{},
	message.TypeAuthPasswordBackendError: message.PayloadAuthPasswordBackendError{},

	message.TypeAuthPubKey:             message.PayloadAuthPubKey{},
	message.TypeAuthPubKeySuccessful:   message.PayloadAuthPubKey{},
	message.TypeAuthPubKeyFailed:       message.PayloadAuthPubKey{},
	message.TypeAuthPubKeyBackendError: message.PayloadAuthPubKeyBackendError{},

	message.TypeGlobalRequestUnknown: message.PayloadGlobalRequestUnknown{},
	message.TypeNewChannel:           message.PayloadNewChannel{},
	message.TypeNewChannelSuccessful: message.PayloadNewChannelSuccessful{},
	message.TypeNewChannelFailed:     message.PayloadNewChannelFailed{},

	message.TypeChannelRequestUnknownType:  message.PayloadChannelRequestUnknownType{},
	message.TypeChannelRequestDecodeFailed: message.PayloadChannelRequestDecodeFailed{},
	message.TypeChannelRequestSetEnv:       message.PayloadChannelRequestSetEnv{},
	message.TypeChannelRequestExec:         message.PayloadChannelRequestExec{},
	message.TypeChannelRequestPty:          message.PayloadChannelRequestPty{},
	message.TypeChannelRequestShell:        message.PayloadChannelRequestShell{},
	message.TypeChannelRequestSignal:       message.PayloadChannelRequestSignal{},
	message.TypeChannelRequestSubsystem:    message.PayloadChannelRequestSubsystem{},
	message.TypeChannelRequestWindow:       message.PayloadChannelRequestWindow{},
	message.TypeIO:                         message.PayloadIO{},
}

func decodeMessage(v decodedMessage) (*message.Message, error) {
	payload, ok := messageTypeMap[v.MessageType]
	if !ok {
		return nil, fmt.Errorf("invalid message type: %d", v.MessageType)
	}

	if payload != nil {
		if err := mapstructure.Decode(v.Payload, &payload); err != nil {
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
