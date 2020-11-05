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

func (e *encoder) sendHeader(header header, storage io.Writer) error {
	data, err := json.Marshal(header)
	if err != nil {
		return err
	}
	_, err = storage.Write(append(data, '\n'))
	if err != nil {
		return err
	}
	return nil
}

func (e *encoder) sendFrame(frame Frame, storage io.Writer) error {
	data, err := frame.MarshalJSON()
	if err != nil {
		return fmt.Errorf("failed to marshal Asciicast frame (%w)", err)
	}
	if _, err = storage.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write Asciicast frame (%w)", err)
	}
	return nil
}

func (e *encoder) Encode(messages <-chan message.Message, storage storage.Writer) error {
	asciicastHeader := header{
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
	const shell = "/bin/sh"
	for {
		msg, ok := <-messages
		if !ok {
			break
		}
		var err error
		startTime, headerWritten, err = e.encodeMessage(startTime, msg, &asciicastHeader, ip, storage, username, headerWritten, shell)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *encoder) encodeMessage(startTime int64, msg message.Message, asciicastHeader *header, ip string, storage storage.Writer, username *string, headerWritten bool, shell string) (int64, bool, error) {
	if startTime == 0 {
		startTime = msg.Timestamp
		asciicastHeader.Timestamp = int(startTime / 1000000000)
	}
	var err error
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
		payload := msg.Payload.(*message.PayloadChannelRequestSetEnv)
		asciicastHeader.Env[payload.Name] = payload.Value
	case message.TypeChannelRequestPty:
		payload := msg.Payload.(*message.PayloadChannelRequestPty)
		asciicastHeader.Width = payload.Columns
		asciicastHeader.Height = payload.Rows
	case message.TypeChannelRequestExec:
		payload := msg.Payload.(*message.PayloadChannelRequestExec)
		startTime, headerWritten, err = e.handleRun(startTime, headerWritten, asciicastHeader, payload.Program, storage)
	case message.TypeChannelRequestShell:
		startTime, headerWritten, err = e.handleRun(startTime, headerWritten, asciicastHeader, shell, storage)
	case message.TypeChannelRequestSubsystem:
		startTime, headerWritten, err = e.handleRun(startTime, headerWritten, asciicastHeader, shell, storage)
	case message.TypeIO:
		startTime, headerWritten, err = e.handleIO(startTime, msg, asciicastHeader, headerWritten, shell, storage)
	}
	if err != nil {
		return startTime, headerWritten, err
	}
	return startTime, headerWritten, nil
}

func (e *encoder) handleRun(startTime int64, headerWritten bool, asciicastHeader *header, program string, storage storage.Writer) (int64, bool, error) {
	if !headerWritten {
		asciicastHeader.Command = program
		if err := e.sendHeader(*asciicastHeader, storage); err != nil {
			return startTime, headerWritten, err
		}
		headerWritten = true
	}
	return startTime, headerWritten, nil
}

func (e *encoder) handleIO(startTime int64, msg message.Message, asciicastHeader *header, headerWritten bool, shell string, storage storage.Writer) (int64, bool, error) {
	if !headerWritten {
		asciicastHeader.Command = shell
		if err := e.sendHeader(*asciicastHeader, storage); err != nil {
			return startTime, headerWritten, err
		}
		headerWritten = true
	}
	payload := msg.Payload.(*message.PayloadIO)
	if payload.Stream == message.StreamStdout ||
		payload.Stream == message.StreamStderr {
		time := float64(msg.Timestamp-startTime) / 1000000000
		frame := Frame{
			Time:      time,
			EventType: eventTypeOutput,
			Data:      string(payload.Data),
		}
		if err := e.sendFrame(frame, storage); err != nil {
			return startTime, headerWritten, err
		}
	}
	return startTime, headerWritten, nil
}
