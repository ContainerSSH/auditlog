package asciinema

import "github.com/containerssh/auditlog/codec"

//goland:noinspection GoUnusedExportedFunction
func NewEncoder() codec.Encoder {
	return &encoder{}
}
