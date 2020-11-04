package asciinema

import "encoding/json"

type Header struct {
	Version   uint              `json:"version"`
	Width     uint              `json:"width"`
	Height    uint              `json:"height"`
	Timestamp int               `json:"timestamp"`
	Command   string            `json:"command"`
	Title     string            `json:"title"`
	Env       map[string]string `json:"env"`
}

type EventType string

//goland:noinspection GoUnusedConst
const (
	EventTypeOutput EventType = "o"
	EventTypeInput  EventType = "i"
)

type Frame struct {
	Time      float64
	EventType EventType
	Data      string
}

func (f *Frame) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{f.Time, f.EventType, f.Data})
}
