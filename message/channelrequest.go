package message

// PayloadChannelRequestUnknownType is a payload signaling that a channel request was not supported.
type PayloadChannelRequestUnknownType struct {
	RequestType string `json:"requestType" yaml:"requestType"`
}

// Equals compares two PayloadChannelRequestUnknownType payloads.
func (p PayloadChannelRequestUnknownType) Equals(other Payload) bool {
	p2, ok := other.(PayloadChannelRequestUnknownType)
	if !ok {
		return false
	}
	return p.RequestType == p2.RequestType
}

// PayloadChannelRequestDecodeFailed is a payload that signals a supported request that the server was unable to decode.
type PayloadChannelRequestDecodeFailed struct {
	RequestType string `json:"requestType" yaml:"requestType"`
	Reason      string `json:"reason" yaml:"reason"`
}

// Equals compares two PayloadChannelRequestDecodeFailed payloads.
func (p PayloadChannelRequestDecodeFailed) Equals(other Payload) bool {
	p2, ok := other.(PayloadChannelRequestDecodeFailed)
	if !ok {
		return false
	}
	return p.RequestType == p2.RequestType && p.Reason == p2.Reason
}

// PayloadChannelRequestSetEnv is a payload signaling the request for an environment variable.
type PayloadChannelRequestSetEnv struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

// Equals compares two PayloadChannelRequestSetEnv payloads.
func (p PayloadChannelRequestSetEnv) Equals(other Payload) bool {
	p2, ok := other.(PayloadChannelRequestSetEnv)
	if !ok {
		return false
	}
	return p.Name == p2.Name && p.Value == p2.Value
}

// PayloadChannelRequestExec is a payload signaling the request to execute a program.
type PayloadChannelRequestExec struct {
	Program string `json:"program" yaml:"program"`
}

// Equals compares two PayloadChannelRequestExec payloads.
func (p PayloadChannelRequestExec) Equals(other Payload) bool {
	p2, ok := other.(PayloadChannelRequestExec)
	if !ok {
		return false
	}
	return p.Program == p2.Program
}

// PayloadChannelRequestPty is a payload signaling the request for an interactive terminal.
type PayloadChannelRequestPty struct {
	Columns uint `json:"columns" yaml:"columns"`
	Rows    uint `json:"rows" yaml:"rows"`
}

// Equals compares two PayloadChannelRequestPty payloads.
func (p PayloadChannelRequestPty) Equals(other Payload) bool {
	p2, ok := other.(PayloadChannelRequestPty)
	if !ok {
		return false
	}
	return p.Columns == p2.Columns && p.Rows == p2.Rows
}

// PayloadChannelRequestShell is a payload signaling a request for a shell.
type PayloadChannelRequestShell struct {
}

// Equals compares two PayloadChannelRequestShell payloads.
func (p PayloadChannelRequestShell) Equals(other Payload) bool {
	_, ok := other.(PayloadChannelRequestShell)
	return ok
}

// PayloadChannelRequestSignal is a payload signaling a signal request to be sent to the currently running program.
type PayloadChannelRequestSignal struct {
	Signal string `json:"signal" yaml:"signal"`
}

// Equals compares two PayloadChannelRequestSignal payloads.
func (p PayloadChannelRequestSignal) Equals(other Payload) bool {
	p2, ok := other.(PayloadChannelRequestSignal)
	if !ok {
		return false
	}
	return p.Signal == p2.Signal
}

// PayloadChannelRequestSubsystem is a payload requesting a well-known subsystem (e.g. sftp)
type PayloadChannelRequestSubsystem struct {
	Subsystem string `json:"subsystem" yaml:"subsystem"`
}

// Equals compares two PayloadChannelRequestSubsystem payloads.
func (p PayloadChannelRequestSubsystem) Equals(other Payload) bool {
	p2, ok := other.(PayloadChannelRequestSubsystem)
	if !ok {
		return false
	}
	return p.Subsystem == p2.Subsystem
}

// PayloadChannelRequestWindow is a payload requesting the change in the terminal window size.
type PayloadChannelRequestWindow struct {
	Columns uint `json:"columns" yaml:"columns"`
	Rows    uint `json:"rows" yaml:"rows"`
}

// Equals compares two PayloadChannelRequestWindow payloads.
func (p PayloadChannelRequestWindow) Equals(other Payload) bool {
	p2, ok := other.(PayloadChannelRequestWindow)
	if !ok {
		return false
	}
	return p.Columns == p2.Columns && p.Rows == p2.Rows
}
