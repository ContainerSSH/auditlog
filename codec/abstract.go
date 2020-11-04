package codec

import (
	"io"

	"github.com/containerssh/auditlog/message"
	"github.com/containerssh/auditlog/storage"
)

type Encoder interface {
	Encode(messages <-chan message.Message, storage storage.Writer) error
	GetMimeType() string
	GetFileExtension() string
}

type Decoder interface {
	Decode(reader io.Reader) (<-chan *message.Message, <-chan error, <-chan bool)
}
