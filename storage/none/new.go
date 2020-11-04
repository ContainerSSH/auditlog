package none

import "github.com/containerssh/auditlog/storage"

func NewStorage() storage.WritableStorage {
	return &Storage{}
}
