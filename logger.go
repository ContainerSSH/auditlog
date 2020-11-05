package auditlog

import (
	"io"
	"net"

	"github.com/containerssh/auditlog/message"
)

type Logger interface {
	OnConnect(connectionID message.ConnectionID, ip net.TCPAddr) (Connection, error)
}

type Connection interface {
	OnDisconnect()

	OnAuthPassword(username string, password []byte)
	OnAuthPasswordSuccess(username string, password []byte)
	OnAuthPasswordFailed(username string, password []byte)
	OnAuthPasswordBackendError(username string, password []byte, reason string)

	OnAuthPubKey(username string, pubKey []byte)
	OnAuthPubKeySuccess(username string, pubKey []byte)
	OnAuthPubKeyFailed(username string, pubKey []byte)
	OnAuthPubKeyBackendError(username string, pubKey []byte, reason string)

	OnGlobalRequestUnknown(requestType string)

	OnNewChannel(channelType string)
	OnNewChannelFailed(channelType string, reason string)
	OnNewChannelSuccess(channelType string, channelID message.ChannelID) Channel
}

type Channel interface {
	OnRequestUnknown(requestType string)
	OnRequestDecodeFailed(requestType string, reason string)

	OnRequestSetEnv(name string, value string)
	OnRequestExec(program string)
	OnRequestPty(columns uint, rows uint)
	OnRequestShell()
	OnRequestSignal(signal string)
	OnRequestSubsystem(subsystem string)
	OnRequestWindow(columns uint, rows uint)

	GetStdinProxy(stdin io.Reader) io.Reader
	GetStdoutProxy(stdout io.Writer) io.Writer
	GetStderrProxy(stderr io.Writer) io.Writer
}
