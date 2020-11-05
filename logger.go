package auditlog

import (
	"io"
	"net"

	"github.com/containerssh/auditlog/message"
)

// Logger is a top level audit logger.
type Logger interface {
	// OnConnect creates an audit log message for a new connection and simultaneously returns a connection object for
	//           connection-specific messages
	OnConnect(connectionID message.ConnectionID, ip net.TCPAddr) (Connection, error)
}

// Connection is an audit logger for a specific connection
type Connection interface {
	// OnDisconnect creates an audit log message for a disconnect event.
	OnDisconnect()

	// OnAuthPassword creates an audit log message for an authentication attempt.
	OnAuthPassword(username string, password []byte)
	// OnAuthPasswordSuccess creates an audit log message for a successful authentication.
	OnAuthPasswordSuccess(username string, password []byte)
	// OnAuthPasswordFailed creates an audit log message for a failed authentication.
	OnAuthPasswordFailed(username string, password []byte)
	// OnAuthPasswordBackendError creates an audit log message for an auth server (backend) error during password
	//                            verification.
	OnAuthPasswordBackendError(username string, password []byte, reason string)

	// OnAuthPubKey creates an audit log message for an authentication attempt with public key.
	OnAuthPubKey(username string, pubKey []byte)
	// OnAuthPubKeySuccess creates an audit log message for a successful public key authentication.
	OnAuthPubKeySuccess(username string, pubKey []byte)
	// OnAuthPubKeyFailed creates an audit log message for a failed public key authentication.
	OnAuthPubKeyFailed(username string, pubKey []byte)
	// OnAuthPubKeyBackendError creates an audit log message for a failure while talking to the auth server (backend)
	//                          during public key authentication.
	OnAuthPubKeyBackendError(username string, pubKey []byte, reason string)

	// OnGlobalRequestUnknown creates an audit log message for a global request that is not supported.
	OnGlobalRequestUnknown(requestType string)

	// OnNewChannel creates an audit log message for a new channel request.
	OnNewChannel(channelType string)
	// OnNewChannelFailed creates an audit log message for a failure in requesting a new channel.
	OnNewChannelFailed(channelType string, reason string)
	// OnNewChannelSuccess creates an audit log message for successfully requesting a new channel and returns a
	//                     channel-specific audit logger.
	OnNewChannelSuccess(channelType string, channelID message.ChannelID) Channel
}

// Channel is an audit logger for one specific hannel
type Channel interface {
	// OnRequestUnknown creates an audit log message for a channel request that is not supported.
	OnRequestUnknown(requestType string)
	// OnRequestDecodeFailed creates an audit log message for a channel request that is supported but could not be
	//                       decoded.
	OnRequestDecodeFailed(requestType string, reason string)

	// OnRequestSetEnv creates an audit log message for a channel request to set an environment variable.
	OnRequestSetEnv(name string, value string)
	// OnRequestExec creates an audit log message for a channel request to execute a program.
	OnRequestExec(program string)
	// OnRequestPty creates an audit log message for a channel request to create an interactive terminal.
	OnRequestPty(columns uint, rows uint)
	// OnRequestExec creates an audit log message for a channel request to execute a shell.
	OnRequestShell()
	// OnRequestExec creates an audit log message for a channel request to send a signal to the currently running
	//               program.
	OnRequestSignal(signal string)
	// OnRequestExec creates an audit log message for a channel request to execute a well-known subsystem (e.g. SFTP)
	OnRequestSubsystem(subsystem string)
	// OnRequestWindow creates an audit log message for a channel request to resize the current window.
	OnRequestWindow(columns uint, rows uint)

	// GetStdinProxy creates an intercepting audit log reader proxy for the standard input.
	GetStdinProxy(stdin io.Reader) io.Reader
	// GetStdinProxy creates an intercepting audit log writer proxy for the standard output.
	GetStdoutProxy(stdout io.Writer) io.Writer
	// GetStdinProxy creates an intercepting audit log writer proxy for the standard error.
	GetStderrProxy(stderr io.Writer) io.Writer
}
