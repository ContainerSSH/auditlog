package asciinema

import "github.com/containerssh/auditlog/codec"

// NewEncoder Creates an encoder that writes in the Asciicast v2 format
//            (see https://github.com/asciinema/asciinema/blob/develop/doc/asciicast-v2.md)
func NewEncoder() codec.Encoder {
	return &encoder{}
}
