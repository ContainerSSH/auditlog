package codec

import (
	"github.com/containerssh/containerssh-auditlog-go/message"
	"io"
)

type Encoder interface {
	Encode(messages <-chan message.Message, storage StorageWriter) error
	GetMimeType() string
	GetFileExtension() string
}

type Decoder interface {
	Decode(reader io.Reader) (<-chan *message.Message, <-chan error, <-chan bool)
}

// The StorageWriter is a regular WriteCloser with an added function to set the connection metadata for indexing.
type StorageWriter interface {
	io.WriteCloser

	// Set metadata for the audit log. Can be called multiple times.
	//
	// startTime is the time when the connection started in unix timestamp
	// sourceIp  is the IP address the user connected from
	// username  is the username the user entered. The first time this method is called the username will be nil,
	//           may be called subsequently is the user authenticated.
	SetMetadata(startTime int64, sourceIp string, username *string)
}