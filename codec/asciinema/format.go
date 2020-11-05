package asciinema

import "encoding/json"

type header struct {
	Version   uint              `json:"version"`
	Width     uint              `json:"width"`
	Height    uint              `json:"height"`
	Timestamp int               `json:"timestamp"`
	Command   string            `json:"command"`
	Title     string            `json:"title"`
	Env       map[string]string `json:"env"`
}

type eventType string

//goland:noinspection GoUnusedConst
const (
	eventTypeOutput eventType = "o"
	eventTypeInput  eventType = "i"
)

type Frame struct {
	Time      float64
	EventType eventType
	Data      string
}

func (f *Frame) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{f.Time, f.EventType, f.Data})
}
