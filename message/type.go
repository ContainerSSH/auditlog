package message

// Type Is the ID for the message type describing which payload is in the message.
type Type int32

const (
	// TypeConnect describes a message that is sent when the user connects on a TCP level.
	TypeConnect Type = 0
	// TypeDisconnect describes a message that is sent when the user disconnects on a TCP level.
	TypeDisconnect Type = 1
	// TypeAuthPassword describes a message that is sent when the user submits a username and password.
	TypeAuthPassword Type = 100
	// TypeAuthPasswordSuccessful describes a message that is sent when the submitted username and password were valid.
	TypeAuthPasswordSuccessful Type = 101
	// TypeAuthPasswordFailed describes a message that is sent when the submitted username and password were invalid.
	TypeAuthPasswordFailed Type = 102
	// TypeAuthPasswordBackendError describes a message that is sent when the auth server failed to respond to a request
	//                              with username and password
	TypeAuthPasswordBackendError Type = 103
	// TypeAuthPubKey describes a message that is sent when the user submits a username and public key.
	TypeAuthPubKey Type = 104
	// TypeAuthPubKeySuccessful describes a message that is sent when the submitted username and public key were invalid.
	TypeAuthPubKeySuccessful Type = 105
	// TypeAuthPubKeyFailed describes a message that is sent when the submitted username and public key were invalid.
	TypeAuthPubKeyFailed Type = 106
	// TypeAuthPubKeyBackendError describes a message that is sent when the auth server failed to respond with username
	//                            and password.
	TypeAuthPubKeyBackendError Type = 107
	// TypeGlobalRequestUnknown describes a message when a global (non-channel) request was sent that was not recognized.
	TypeGlobalRequestUnknown Type = 200
	// TypeNewChannel describes a message that indicates a new channel request
	TypeNewChannel Type = 300
	// TypeNewChalleSuccessful describes a message when the new channel request was successful
	TypeNewChannelSuccessful Type = 301
	// TypeNewChannelFailed describes a message when the channel request failed for the reason indicated
	TypeNewChannelFailed Type = 302
	// TypeChannelRequestUnknownType describes an in-channel request from the user that is not supported
	TypeChannelRequestUnknownType Type = 400
	// TypeChannelRequestDecodeFailed describes an in-channel request from the user that is supported but the payload
	//                                could not be decoded.
	TypeChannelRequestDecodeFailed Type = 401
	// TypeChannelRequestSetEnv describes an in-channel request to set an environment variable
	TypeChannelRequestSetEnv Type = 402
	// TypeChannelRequestSetEnv describes an in-channel request to run a program
	TypeChannelRequestExec Type = 403
	// TypeChannelRequestSetEnv describes an in-channel request to create an interactive terminal
	TypeChannelRequestPty Type = 404
	// TypeChannelRequestSetEnv describes an in-channel request to start a shell
	TypeChannelRequestShell Type = 405
	// TypeChannelRequestSetEnv describes an in-channel request to send a signal to the currently running program
	TypeChannelRequestSignal Type = 406
	// TypeChannelRequestSetEnv describes an in-channel request to start a well-known subsystem (e.g. SFTP)
	TypeChannelRequestSubsystem Type = 407
	// TypeChannelRequestSetEnv describes an in-channel request to resize the current interactive terminal
	TypeChannelRequestWindow Type = 408
	// TypeChannelRequestSetEnv describes the data transfered to and from the currently running program on the terminal.
	TypeIO Type = 500
)

var typeToName = map[Type]string{
	TypeConnect:    "connect",
	TypeDisconnect: "disconnect",

	TypeAuthPassword:             "auth_password",
	TypeAuthPasswordSuccessful:   "auth_password_successful",
	TypeAuthPasswordFailed:       "auth_password_failed",
	TypeAuthPasswordBackendError: "auth_password_backend_error",

	TypeAuthPubKey:             "auth_pubkey",
	TypeAuthPubKeySuccessful:   "auth_pubkey_successful",
	TypeAuthPubKeyFailed:       "auth_pubkey_failed",
	TypeAuthPubKeyBackendError: "auth_pubkey_backend_error",

	TypeGlobalRequestUnknown: "global_request_unknown",
	TypeNewChannel:           "new_channel",
	TypeNewChannelSuccessful: "new_channel_successful",
	TypeNewChannelFailed:     "new_channel_failed",

	TypeChannelRequestUnknownType:  "channel_request_unknown",
	TypeChannelRequestDecodeFailed: "channel_request_decode_failed",
	TypeChannelRequestSetEnv:       "setenv",
	TypeChannelRequestExec:         "exec",
	TypeChannelRequestPty:          "pty",
	TypeChannelRequestShell:        "shell",
	TypeChannelRequestSignal:       "signal",
	TypeChannelRequestSubsystem:    "subsystem",
	TypeChannelRequestWindow:       "window",

	TypeIO: "io",
}

func (messageType Type) ToName() string {
	if val, ok := typeToName[messageType]; ok {
		return val
	}
	return "invalid"
}
