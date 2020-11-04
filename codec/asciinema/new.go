package asciinema

import "github.com/containerssh/containerssh-auditlog-go/codec"

//goland:noinspection GoUnusedExportedFunction
func NewEncoder() codec.Encoder {
	return &encoder{}
}
