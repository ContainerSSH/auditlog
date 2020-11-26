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

	"github.com/containerssh/geoip"
	"github.com/containerssh/log"
)

type loggerImplementation struct {
	intercept   InterceptConfig
	encoder     codec.Encoder
	storage     storage.WritableStorage
	logger      log.Logger
	wg          *sync.WaitGroup
	geoIPLookup geoip.LookupProvider
}

type loggerConnection struct {
	l *loggerImplementation

	ip             net.TCPAddr
	messageChannel chan message.Message
	connectionID   message.ConnectionID
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
			Country:    l.geoIPLookup.Lookup(ip.IP),
		},
		ChannelID: nil,
	}
	return conn, nil
}

func (c *loggerConnection) OnDisconnect() {
	c.messageChannel <- message.Message{
		ConnectionID: c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeDisconnect,
		Payload:      nil,
		ChannelID:    nil,
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
		ChannelID: nil,
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
		ChannelID: nil,
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
		ChannelID: nil,
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
		ChannelID: nil,
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
		ChannelID: nil,
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
		ChannelID: nil,
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
		ChannelID: nil,
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
		ChannelID: nil,
	}
}

func (c *loggerConnection) OnHandshakeFailed(reason string) {
	c.messageChannel <- message.Message{
		ConnectionID: c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeHandshakeFailed,
		Payload: message.PayloadHandshakeFailed{
			Reason: reason,
		},
		ChannelID: nil,
	}
}

func (c *loggerConnection) OnHandshakeSuccessful(username string) {
	c.messageChannel <- message.Message{
		ConnectionID: c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeHandshakeSuccessful,
		Payload: message.PayloadHandshakeSuccessful{
			Username: username,
		},
		ChannelID: nil,
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
		ChannelID: nil,
	}
}

func (c *loggerConnection) OnNewChannel(channelID message.ChannelID, channelType string) {
	c.messageChannel <- message.Message{
		ConnectionID: c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeNewChannel,
		Payload: message.PayloadNewChannel{
			ChannelType: channelType,
		},
		ChannelID: channelID,
	}
}

func (c *loggerConnection) OnNewChannelFailed(channelID message.ChannelID, channelType string, reason string) {
	c.messageChannel <- message.Message{
		ConnectionID: c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeNewChannelFailed,
		Payload: message.PayloadNewChannelFailed{
			ChannelType: channelType,
			Reason:      reason,
		},
		ChannelID: channelID,
	}
}

func (c *loggerConnection) OnNewChannelSuccess(channelID message.ChannelID, channelType string) Channel {
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

func (l *loggerChannel) OnRequestUnknown(requestID uint64, requestType string, payload []byte) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeChannelRequestUnknownType,
		Payload: message.PayloadChannelRequestUnknownType{
			RequestID:   requestID,
			RequestType: requestType,
			Payload:     payload,
		},
		ChannelID: l.channelID,
	}
}

func (l *loggerChannel) OnRequestDecodeFailed(requestID uint64, requestType string, payload []byte, reason string) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeChannelRequestDecodeFailed,
		Payload: message.PayloadChannelRequestDecodeFailed{
			RequestID:   requestID,
			RequestType: requestType,
			Payload:     payload,
			Reason:      reason,
		},
		ChannelID: l.channelID,
	}
}

func (l *loggerChannel) OnRequestSetEnv(requestID uint64, name string, value string) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeChannelRequestSetEnv,
		Payload: message.PayloadChannelRequestSetEnv{
			RequestID: requestID,
			Name:      name,
			Value:     value,
		},
		ChannelID: l.channelID,
	}
}

func (l *loggerChannel) OnRequestExec(requestID uint64, program string) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeChannelRequestExec,
		Payload: message.PayloadChannelRequestExec{
			RequestID: requestID,
			Program:   program,
		},
		ChannelID: l.channelID,
	}
}

func (l *loggerChannel) OnRequestPty(requestID uint64, term string, columns uint32, rows uint32, width uint32, height uint32, modelist []byte) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeChannelRequestPty,
		Payload: message.PayloadChannelRequestPty{
			RequestID: requestID,
			Term:      term,
			Columns:   columns,
			Rows:      rows,
			Width:     width,
			Height:    height,
			ModeList:  modelist,
		},
		ChannelID: l.channelID,
	}
}

func (l *loggerChannel) OnRequestShell(requestID uint64) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeChannelRequestShell,
		Payload: message.PayloadChannelRequestShell{
			RequestID: requestID,
		},
		ChannelID: l.channelID,
	}
}

func (l *loggerChannel) OnRequestSignal(requestID uint64, signal string) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeChannelRequestSignal,
		Payload: message.PayloadChannelRequestSignal{
			RequestID: requestID,
			Signal:    signal,
		},
		ChannelID: l.channelID,
	}
}

func (l *loggerChannel) OnRequestSubsystem(requestID uint64, subsystem string) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeChannelRequestSubsystem,
		Payload: message.PayloadChannelRequestSubsystem{
			RequestID: requestID,
			Subsystem: subsystem,
		},
		ChannelID: l.channelID,
	}
}

func (l *loggerChannel) OnRequestWindow(requestID uint64, columns uint32, rows uint32, width uint32, height uint32) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeChannelRequestWindow,
		Payload: message.PayloadChannelRequestWindow{
			RequestID: requestID,
			Columns:   columns,
			Rows:      rows,
			Width:     width,
			Height:    height,
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

func (l *loggerChannel) OnRequestFailed(requestID uint64, reason error) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeRequestFailed,
		Payload: message.PayloadRequestFailed{
			RequestID: requestID,
			Reason:    reason.Error(),
		},
		ChannelID: l.channelID,
	}
}

func (l *loggerChannel) OnExit(exitStatus uint32) {
	l.c.messageChannel <- message.Message{
		ConnectionID: l.c.connectionID,
		Timestamp:    time.Now().UnixNano(),
		MessageType:  message.TypeExit,
		Payload: message.PayloadExit{
			ExitStatus: exitStatus,
		},
		ChannelID: l.channelID,
	}
}

//endregion
