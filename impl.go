package auditlog

import (
	"context"
	"encoding/hex"
	"io"
	"net"
	"sync"
	"time"

	"github.com/containerssh/auditlog/codec"
	"github.com/containerssh/auditlog/message"
	"github.com/containerssh/auditlog/storage"

	"github.com/containerssh/log"
)

type loggerImplementation struct {
	intercept InterceptConfig
	encoder   codec.Encoder
	storage   storage.WritableStorage
	logger    log.Logger
	wg        *sync.WaitGroup
}

type loggerConnection struct {
	l *loggerImplementation

	ip             net.TCPAddr
	messageChannel chan message.Message
	connectionID   message.ConnectionID
	nextChannelID  message.ChannelID
	lock           *sync.Mutex
}

type loggerChannel struct {
	c *loggerConnection

	channelID message.ChannelID
}

func (l *loggerImplementation) Shutdown(shutdownContext context.Context) {
	l.wg.Wait()
	l.storage.Shutdown(shutdownContext)
}

//region Connection
func (l *loggerImplementation) OnConnect(connectionID message.ConnectionID, ip net.TCPAddr) (Connection, error) {
	name := hex.EncodeToString(connectionID)
	writer, err := l.storage.OpenWriter(name)
	if err != nil {
		return nil, err
	}
	conn := &loggerConnection{
		l:              l,
		ip:             ip,
		connectionID:   connectionID,
		messageChannel: make(chan message.Message),
		lock:           &sync.Mutex{},
		nextChannelID:  0,
	}
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		err := l.encoder.Encode(conn.messageChannel, writer)
		if err != nil {
			l.logger.Emergencye(err)
		}
	}()
	conn.messageChannel <- message.Message{
		ConnectionID: connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeConnect,
		Payload: message.PayloadConnect{
			RemoteAddr: ip.IP.String(),
		},
		ChannelID: -1,
	}
	return conn, nil
}

func (c *loggerConnection) OnDisconnect() {
	c.messageChannel <- message.Message{
		ConnectionID: c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeDisconnect,
		Payload:      nil,
		ChannelID:    -1,
	}
	close(c.messageChannel)
}

func (c *loggerConnection) OnAuthPassword(username string, password []byte) {
	c.messageChannel <- message.Message{
		ConnectionID: c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeAuthPassword,
		Payload: message.PayloadAuthPassword{
			Username: username,
			Password: password,
		},
		ChannelID: -1,
	}
}

func (c *loggerConnection) OnAuthPasswordSuccess(username string, password []byte) {
	c.messageChannel <- message.Message{
		ConnectionID: c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeAuthPasswordSuccessful,
		Payload: message.PayloadAuthPassword{
			Username: username,
			Password: password,
		},
		ChannelID: -1,
	}
}

func (c *loggerConnection) OnAuthPasswordFailed(username string, password []byte) {
	c.messageChannel <- message.Message{
		ConnectionID: c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeAuthPasswordFailed,
		Payload: message.PayloadAuthPassword{
			Username: username,
			Password: password,
		},
		ChannelID: -1,
	}
}

func (c *loggerConnection) OnAuthPasswordBackendError(username string, password []byte, reason string) {
	c.messageChannel <- message.Message{
		ConnectionID: c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeAuthPasswordBackendError,
		Payload: message.PayloadAuthPasswordBackendError{
			Username: username,
			Password: password,
			Reason:   reason,
		},
		ChannelID: -1,
	}
}

func (c *loggerConnection) OnAuthPubKey(username string, pubKey []byte) {
	c.messageChannel <- message.Message{
		ConnectionID: c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeAuthPubKey,
		Payload: message.PayloadAuthPubKey{
			Username: username,
			Key:      pubKey,
		},
		ChannelID: -1,
	}
}

func (c *loggerConnection) OnAuthPubKeySuccess(username string, pubKey []byte) {
	c.messageChannel <- message.Message{
		ConnectionID: c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeAuthPubKeySuccessful,
		Payload: message.PayloadAuthPubKey{
			Username: username,
			Key:      pubKey,
		},
		ChannelID: -1,
	}
}

func (c *loggerConnection) OnAuthPubKeyFailed(username string, pubKey []byte) {
	c.messageChannel <- message.Message{
		ConnectionID: c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeAuthPubKeyFailed,
		Payload: message.PayloadAuthPubKey{
			Username: username,
			Key:      pubKey,
		},
		ChannelID: -1,
	}
}

func (c *loggerConnection) OnAuthPubKeyBackendError(username string, pubKey []byte, reason string) {
	c.messageChannel <- message.Message{
		ConnectionID: c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeAuthPubKeyBackendError,
		Payload: message.PayloadAuthPubKeyBackendError{
			Username: username,
			Key:      pubKey,
			Reason:   reason,
		},
		ChannelID: -1,
	}
}

func (c *loggerConnection) OnGlobalRequestUnknown(requestType string) {
	c.messageChannel <- message.Message{
		ConnectionID: c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeGlobalRequestUnknown,
		Payload: message.PayloadGlobalRequestUnknown{
			RequestType: requestType,
		},
		ChannelID: -1,
	}
}

func (c *loggerConnection) OnNewChannel(channelType string) {
	c.messageChannel <- message.Message{
		ConnectionID: c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeNewChannel,
		Payload: message.PayloadNewChannel{
			ChannelType: channelType,
		},
		ChannelID: -1,
	}
}

func (c *loggerConnection) OnNewChannelFailed(channelType string, reason string) {
	c.messageChannel <- message.Message{
		ConnectionID: c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeNewChannelFailed,
		Payload: message.PayloadNewChannelFailed{
			ChannelType: channelType,
			Reason:      reason,
		},
		ChannelID: -1,
	}
}

func (c *loggerConnection) OnNewChannelSuccess(channelType string) Channel {
	c.lock.Lock()
	channelID := c.nextChannelID
	c.nextChannelID++
	c.lock.Unlock()
	c.messageChannel <- message.Message{
		ConnectionID: c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeNewChannelSuccessful,
		Payload: message.PayloadNewChannelSuccessful{
			ChannelType: channelType,
		},
		ChannelID: channelID,
	}
	return &loggerChannel{
		c:         c,
		channelID: channelID,
	}
}

//endregion

//region Channel
func (l *loggerChannel) OnRequestUnknown(requestType string) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeChannelRequestUnknownType,
		Payload: message.PayloadChannelRequestUnknownType{
			RequestType: requestType,
		},
		ChannelID: l.channelID,
	}
}

func (l *loggerChannel) OnRequestDecodeFailed(requestType string, reason string) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeChannelRequestDecodeFailed,
		Payload: message.PayloadChannelRequestDecodeFailed{
			RequestType: requestType,
			Reason:      reason,
		},
		ChannelID: l.channelID,
	}
}

func (l *loggerChannel) OnRequestSetEnv(name string, value string) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeChannelRequestSetEnv,
		Payload: message.PayloadChannelRequestSetEnv{
			Name:  name,
			Value: value,
		},
		ChannelID: l.channelID,
	}
}

func (l *loggerChannel) OnRequestExec(program string) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeChannelRequestExec,
		Payload: message.PayloadChannelRequestExec{
			Program: program,
		},
		ChannelID: l.channelID,
	}
}

func (l *loggerChannel) OnRequestPty(columns uint, rows uint) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeChannelRequestPty,
		Payload: message.PayloadChannelRequestPty{
			Columns: columns,
			Rows:    rows,
		},
		ChannelID: l.channelID,
	}
}

func (l *loggerChannel) OnRequestShell() {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeChannelRequestShell,
		Payload:      message.PayloadChannelRequestShell{},
		ChannelID:    l.channelID,
	}
}

func (l *loggerChannel) OnRequestSignal(signal string) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeChannelRequestSignal,
		Payload: message.PayloadChannelRequestSignal{
			Signal: signal,
		},
		ChannelID: l.channelID,
	}
}

func (l *loggerChannel) OnRequestSubsystem(subsystem string) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeChannelRequestSubsystem,
		Payload: message.PayloadChannelRequestSubsystem{
			Subsystem: subsystem,
		},
		ChannelID: l.channelID,
	}
}

func (l *loggerChannel) OnRequestWindow(columns uint, rows uint) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeChannelRequestWindow,
		Payload: message.PayloadChannelRequestWindow{
			Columns: columns,
			Rows:    rows,
		},
		ChannelID: l.channelID,
	}
}

func (l *loggerChannel) io(stream message.Stream, data []byte) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeIO,
		Payload: message.PayloadIO{
			Stream: stream,
			Data:   data,
		},
		ChannelID: l.channelID,
	}
}

func (l *loggerChannel) GetStdinProxy(stdin io.Reader) io.Reader {
	if !l.c.l.intercept.Stdin {
		return stdin
	}
	return &interceptingReader{
		backend: stdin,
		stream:  message.StreamStdin,
		channel: l,
	}
}

func (l *loggerChannel) GetStdoutProxy(stdout io.Writer) io.Writer {
	if !l.c.l.intercept.Stdout {
		return stdout
	}
	return &interceptingWriter{
		backend: stdout,
		stream:  message.StreamStdout,
		channel: l,
	}
}

func (l *loggerChannel) GetStderrProxy(stderr io.Writer) io.Writer {
	if !l.c.l.intercept.Stdout {
		return stderr
	}
	return &interceptingWriter{
		backend: stderr,
		stream:  message.StreamStderr,
		channel: l,
	}
}

//endregion
