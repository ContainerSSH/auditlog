package storage

import "io"

type Entry struct {
	Name     string
	Metadata map[string]string
}

type ReadWriteStorage interface {
	ReadableStorage
	WritableStorage
}

type WritableStorage interface {
	OpenWriter(name string) (Writer, error)
}

type ReadableStorage interface {
	OpenReader(name string) (io.ReadCloser, error)
	List() (<-chan Entry, <-chan error)
}

// The Writer is a regular WriteCloser with an added function to set the connection metadata for indexing.
type Writer interface {
	io.WriteCloser

	// Set metadata for the audit log. Will be called multiple times, once when user connects and once when the user
	// authenticates.
	//
	// startTime is the time when the connection started in unix timestamp
	// sourceIp  is the IP address the user connected from
	// username  is the username the user entered. The first time this method is called the username will be nil,
	//           may be called subsequently is the user authenticated.
	SetMetadata(startTime int64, sourceIP string, username *string)
}
