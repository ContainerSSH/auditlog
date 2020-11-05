package codec

import (
	"io"

	"github.com/containerssh/auditlog/message"
	"github.com/containerssh/auditlog/storage"
)

// Encoder is a module that is responsible for receiving audit log messages and writing them to a writer.
type Encoder interface {
	Encode(messages <-chan message.Message, storage storage.Writer) error
	GetMimeType() string
	GetFileExtension() string
}

// Decoder is a module that is resonsible for decoding a binary data stream into audit log messages.
type Decoder interface {
	Decode(reader io.Reader) (<-chan message.Message, <-chan error, <-chan bool)
}
