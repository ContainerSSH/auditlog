package asciinema

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/containerssh/auditlog/message"
	"github.com/containerssh/auditlog/storage"
)

type encoder struct {
}

func (e *encoder) GetMimeType() string {
	return "application/x-asciicast"
}

func (e *encoder) GetFileExtension() string {
	return ".cast"
}

func (e *encoder) sendHeader(header Header, storage io.Writer) error {
	data, err := json.Marshal(header)
	if err != nil {
		return err
	} else {
		_, err = storage.Write(append(data, '\n'))
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *encoder) sendFrame(frame Frame, storage io.Writer) error {
	data, err := frame.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal Asciicast frame (%v)", err)
	} else {
		if _, err = storage.Write(append(data, '\n')); err != nil {
			return fmt.Errorf("failed to write Asciicast frame (%v)", err)
		}
	}
	return nil
}

func (e *encoder) Encode(messages <-chan message.Message, storage storage.Writer) error {
	asciicastHeader := Header{
		Version:   2,
		Width:     80,
		Height:    25,
		Timestamp: 0,
		Command:   "",
		Title:     "",
		Env:       map[string]string{},
	}
	startTime := int64(0)
	headerWritten := false
	var ip = ""
	var username *string
	for {
		msg, ok := <-messages
		if !ok {
			break
		}
		if startTime == 0 {
			startTime = msg.Timestamp
			asciicastHeader.Timestamp = int(startTime / 1000000000)
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
		case message.TypeChannelRequestSetEnv:
			if headerWritten {
				break
			}
			payload := msg.Payload.(*message.PayloadChannelRequestSetEnv)
			asciicastHeader.Env[payload.Name] = payload.Value
		case message.TypeChannelRequestPty:
			if headerWritten {
				break
			}
			payload := msg.Payload.(*message.PayloadChannelRequestPty)
			asciicastHeader.Width = payload.Columns
			asciicastHeader.Height = payload.Rows
		case message.TypeChannelRequestExec:
			if headerWritten {
				break
			}
			payload := msg.Payload.(*message.PayloadChannelRequestExec)
			asciicastHeader.Command = payload.Program
			if err := e.sendHeader(asciicastHeader, storage); err != nil {
				return err
			}
			headerWritten = true
		case message.TypeChannelRequestShell:
			if headerWritten {
				break
			}
			asciicastHeader.Command = "/bin/sh"
			if err := e.sendHeader(asciicastHeader, storage); err != nil {
				return err
			}
			headerWritten = true
		case message.TypeChannelRequestSubsystem:
			//Fallback
			if headerWritten {
				break
			}
			asciicastHeader.Command = "/bin/sh"
			if err := e.sendHeader(asciicastHeader, storage); err != nil {
				return err
			}
			headerWritten = true
		case message.TypeIO:
			if !headerWritten {
				asciicastHeader.Command = "/bin/sh"
				if err := e.sendHeader(asciicastHeader, storage); err != nil {
					return err
				}
				headerWritten = true
			}
			payload := msg.Payload.(*message.PayloadIO)
			if payload.Stream == message.StreamStdout ||
				payload.Stream == message.StreamStderr {
				time := float64(msg.Timestamp-startTime) / 1000000000
				frame := Frame{
					Time:      time,
					EventType: EventTypeOutput,
					Data:      string(payload.Data),
				}
				if err := e.sendFrame(frame, storage); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
