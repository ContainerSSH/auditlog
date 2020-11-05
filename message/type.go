package message

type Type int32

const (
	TypeConnect                    Type = 0
	TypeDisconnect                 Type = 1
	TypeAuthPassword               Type = 100
	TypeAuthPasswordSuccessful     Type = 101
	TypeAuthPasswordFailed         Type = 102
	TypeAuthPasswordBackendError   Type = 103
	TypeAuthPubKey                 Type = 104
	TypeAuthPubKeySuccessful       Type = 105
	TypeAuthPubKeyFailed           Type = 106
	TypeAuthPubKeyBackendError     Type = 107
	TypeGlobalRequestUnknown       Type = 200
	TypeNewChannel                 Type = 300
	TypeNewChannelSuccessful       Type = 301
	TypeNewChannelFailed           Type = 302
	TypeChannelRequestUnknownType  Type = 400
	TypeChannelRequestDecodeFailed Type = 401
	TypeChannelRequestSetEnv       Type = 402
	TypeChannelRequestExec         Type = 403
	TypeChannelRequestPty          Type = 404
	TypeChannelRequestShell        Type = 405
	TypeChannelRequestSignal       Type = 406
	TypeChannelRequestSubsystem    Type = 407
	TypeChannelRequestWindow       Type = 408
	TypeIO                         Type = 500
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
