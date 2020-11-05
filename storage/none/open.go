package none

import "github.com/containerssh/auditlog/storage"

func (s nopStorage) OpenWriter(_ string) (storage.Writer, error) {
	return &nullWriteCloser{}, nil
}
