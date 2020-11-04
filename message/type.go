package message

type MessageType int32

const (
	TypeConnect                    MessageType = 0
	TypeDisconnect                 MessageType = 1
	TypeAuthPassword               MessageType = 100
	TypeAuthPasswordSuccessful     MessageType = 101
	TypeAuthPasswordFailed         MessageType = 102
	TypeAuthPasswordBackendError   MessageType = 103
	TypeAuthPubKey                 MessageType = 104
	TypeAuthPubKeySuccessful       MessageType = 105
	TypeAuthPubKeyFailed           MessageType = 106
	TypeAuthPubKeyBackendError     MessageType = 107
	TypeGlobalRequestUnknown       MessageType = 200
	TypeNewChannel                 MessageType = 300
	TypeNewChannelSuccessful       MessageType = 301
	TypeNewChannelFailed           MessageType = 302
	TypeChannelRequestUnknownType  MessageType = 400
	TypeChannelRequestDecodeFailed MessageType = 401
	TypeChannelRequestSetEnv       MessageType = 402
	TypeChannelRequestExec         MessageType = 403
	TypeChannelRequestPty          MessageType = 404
	TypeChannelRequestShell        MessageType = 405
	TypeChannelRequestSignal       MessageType = 406
	TypeChannelRequestSubsystem    MessageType = 407
	TypeChannelRequestWindow       MessageType = 408
	TypeIO                         MessageType = 500
)

func (messageType MessageType) ToName() string {
	switch messageType {

	case TypeConnect:
		return "auth_connect"
	case TypeDisconnect:
		return "disconnect"

	case TypeAuthPassword:
		return "auth_password"
	case TypeAuthPasswordSuccessful:
		return "auth_password_successful"
	case TypeAuthPasswordFailed:
		return "auth_password_failed"
	case TypeAuthPasswordBackendError:
		return "auth_password_backend_error"

	case TypeAuthPubKey:
		return "auth_pubkey"
	case TypeAuthPubKeySuccessful:
		return "auth_pubkey_successful"
	case TypeAuthPubKeyFailed:
		return "auth_pubkey_failed"
	case TypeAuthPubKeyBackendError:
		return "auth_pubkey_backend_error"

	case TypeGlobalRequestUnknown:
		return "global_request_unknown"
	case TypeNewChannel:
		return "new_channel"
	case TypeNewChannelSuccessful:
		return "new_channel_successful"
	case TypeNewChannelFailed:
		return "new_channel_failed"

	case TypeChannelRequestUnknownType:
		return "channel_request_unknown"
	case TypeChannelRequestDecodeFailed:
		return "channel_request_decode_failed"
	case TypeChannelRequestSetEnv:
		return "setenv"
	case TypeChannelRequestExec:
		return "exec"
	case TypeChannelRequestPty:
		return "pty"
	case TypeChannelRequestShell:
		return "shell"
	case TypeChannelRequestSignal:
		return "signal"
	case TypeChannelRequestSubsystem:
		return "subsystem"
	case TypeChannelRequestWindow:
		return "window"
	case TypeIO:
		return "io"
	default:
		return "invalid"
	}
}
