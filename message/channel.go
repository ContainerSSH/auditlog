package message

type PayloadNewChannel struct {
	ChannelType string `json:"channelType" yaml:"channelType"`
}

func (p PayloadNewChannel) Equals(other Payload) bool {
	p2, ok := other.(PayloadNewChannel)
	if !ok {
		return false
	}
	return p.ChannelType == p2.ChannelType
}

type PayloadNewChannelFailed struct {
	ChannelType string `json:"channelType" yaml:"channelType"`
	Reason      string `json:"reason" yaml:"reason"`
}

func (p PayloadNewChannelFailed) Equals(other Payload) bool {
	p2, ok := other.(PayloadNewChannelFailed)
	if !ok {
		return false
	}
	return p.ChannelType == p2.ChannelType && p.Reason == p2.Reason
}

type PayloadNewChannelSuccessful struct {
	ChannelType string `json:"channelType" yaml:"channelType"`
}

func (p PayloadNewChannelSuccessful) Equals(other Payload) bool {
	p2, ok := other.(PayloadNewChannelSuccessful)
	if !ok {
		return false
	}
	return p.ChannelType == p2.ChannelType
}
