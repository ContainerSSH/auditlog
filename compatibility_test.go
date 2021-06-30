package auditlog_test

import (
	"testing"

	"github.com/containerssh/auditlog/codec/binary"
	"github.com/containerssh/auditlog/message"
	"github.com/containerssh/auditlog/storage"
	"github.com/containerssh/auditlog/storage/file"
	"github.com/containerssh/log"
	"github.com/stretchr/testify/assert"
)

func TestDecodingOldAuditLogs(t *testing.T) {
	logger := log.NewTestLogger(t)
	fileStorage, err := file.NewStorage(file.Config{
		Directory: "./testdata/",
	}, logger)
	if err != nil {
		t.Fatal(err)
	}
	logChan, errChan := fileStorage.List()
	loop:
	for {
		var entry storage.Entry
		var ok bool
		var err error
		select {
		case entry, ok = <-logChan:
			if !ok {
				break loop
			}
		case err, ok = <-errChan:
			if !ok {
				break loop
			}
			t.Fatal(err)
		}

		name := entry.Name
		t.Run(name, func(t *testing.T) {
			testDecodeOldLog(t, fileStorage, name)
		})
	}
}

func testDecodeOldLog(t *testing.T, fileStorage storage.ReadWriteStorage, name string) {
	reader, err := fileStorage.OpenReader(name)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = reader.Close()
	}()
	decoder := binary.NewDecoder()
	messageChannel, errors := decoder.Decode(reader)
	var types []message.Type
	loop:
	for {
		select {
		case msg, ok := <-messageChannel:
			if !ok {
				break loop
			}
			types = append(types, msg.MessageType)
		case err, ok := <-errors:
			if !ok {
				break loop
			}
			t.Fatal(err)
		}
	}
	assert.Equal(t, []message.Type{
		message.TypeConnect,
		message.TypeAuthPassword,
		message.TypeAuthPasswordSuccessful,
		message.TypeHandshakeSuccessful,
		message.TypeNewChannelSuccessful,
		message.TypeChannelRequestPty,
		message.TypeChannelRequestShell,
		message.TypeWriteClose,
		message.TypeExit,
		message.TypeClose,
		message.TypeDisconnect,
	}, types)
}