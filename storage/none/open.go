package none

import "github.com/containerssh/auditlog/storage"

func (s Storage) OpenWriter(_ string) (storage.Writer, error) {
	return &nullWriteCloser{}, nil
}
