package message

type PayloadChannelRequestUnknownType struct {
	RequestType string `json:"requestType" yaml:"requestType"`
}

func (p PayloadChannelRequestUnknownType) Equals(other Payload) bool {
	p2, ok := other.(PayloadChannelRequestUnknownType)
	if !ok {
		return false
	}
	return p.RequestType == p2.RequestType
}

type PayloadChannelRequestDecodeFailed struct {
	RequestType string `json:"requestType" yaml:"requestType"`
	Reason      string `json:"reason" yaml:"reason"`
}

func (p PayloadChannelRequestDecodeFailed) Equals(other Payload) bool {
	p2, ok := other.(PayloadChannelRequestDecodeFailed)
	if !ok {
		return false
	}
	return p.RequestType == p2.RequestType && p.Reason == p2.Reason
}

type PayloadChannelRequestFailed struct {
	RequestType string `json:"requestType" yaml:"requestType"`
	Reason      string `json:"reason" yaml:"reason"`
}

func (p PayloadChannelRequestFailed) Equals(other Payload) bool {
	p2, ok := other.(PayloadChannelRequestFailed)
	if !ok {
		return false
	}
	return p.RequestType == p2.RequestType && p.Reason == p2.Reason
}

type PayloadChannelRequestSetEnv struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

func (p PayloadChannelRequestSetEnv) Equals(other Payload) bool {
	p2, ok := other.(PayloadChannelRequestSetEnv)
	if !ok {
		return false
	}
	return p.Name == p2.Name && p.Value == p2.Value
}

type PayloadChannelRequestExec struct {
	Program string `json:"program" yaml:"program"`
}

func (p PayloadChannelRequestExec) Equals(other Payload) bool {
	p2, ok := other.(PayloadChannelRequestExec)
	if !ok {
		return false
	}
	return p.Program == p2.Program
}

type PayloadChannelRequestPty struct {
	Columns uint `json:"columns" yaml:"columns"`
	Rows    uint `json:"rows" yaml:"rows"`
}

func (p PayloadChannelRequestPty) Equals(other Payload) bool {
	p2, ok := other.(PayloadChannelRequestPty)
	if !ok {
		return false
	}
	return p.Columns == p2.Columns && p.Rows == p2.Rows
}

type PayloadChannelRequestShell struct {
}

func (p PayloadChannelRequestShell) Equals(other Payload) bool {
	_, ok := other.(PayloadChannelRequestShell)
	return ok
}

type PayloadChannelRequestSignal struct {
	Signal string `json:"signal" yaml:"signal"`
}

func (p PayloadChannelRequestSignal) Equals(other Payload) bool {
	p2, ok := other.(PayloadChannelRequestSignal)
	if !ok {
		return false
	}
	return p.Signal == p2.Signal
}

type PayloadChannelRequestSubsystem struct {
	Subsystem string `json:"subsystem" yaml:"subsystem"`
}

func (p PayloadChannelRequestSubsystem) Equals(other Payload) bool {
	p2, ok := other.(PayloadChannelRequestSubsystem)
	if !ok {
		return false
	}
	return p.Subsystem == p2.Subsystem
}

type PayloadChannelRequestWindow struct {
	Columns uint `json:"columns" yaml:"columns"`
	Rows    uint `json:"rows" yaml:"rows"`
}

func (p PayloadChannelRequestWindow) Equals(other Payload) bool {
	p2, ok := other.(PayloadChannelRequestWindow)
	if !ok {
		return false
	}
	return p.Columns == p2.Columns && p.Rows == p2.Rows
}
