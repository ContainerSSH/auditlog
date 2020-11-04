package codec

import "io"

func NewStorageWriterProxy(backend io.WriteCloser) StorageWriter {
	return &storageWriterProxy{backend: backend}
}

type storageWriterProxy struct {
	backend io.WriteCloser
}

func (s *storageWriterProxy) Write(p []byte) (n int, err error) {
	return s.backend.Write(p)
}

func (s *storageWriterProxy) Close() error {
	return s.backend.Close()
}

func (s *storageWriterProxy) SetMetadata(_ int64, _ string, _ *string) {
	// No metadata storage
}
