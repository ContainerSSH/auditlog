package auditlog

import (
	"fmt"

	"github.com/containerssh/auditlog/codec"
	"github.com/containerssh/auditlog/codec/asciinema"
	"github.com/containerssh/auditlog/codec/binary"
	noneCodec "github.com/containerssh/auditlog/codec/none"
	"github.com/containerssh/auditlog/storage"
	"github.com/containerssh/auditlog/storage/file"
	noneStorage "github.com/containerssh/auditlog/storage/none"
	"github.com/containerssh/auditlog/storage/s3"

	"github.com/containerssh/log"
)

// New Creates a new audit logging pipeline based on the provided configuration.
//goland:noinspection GoUnusedExportedFunction
func New(config Config, logger log.Logger) (Logger, error) {
	encoder, err := NewEncoder(config.Format, logger)
	if err != nil {
		return nil, err
	}

	st, err := NewStorage(config, logger)
	if err != nil {
		return nil, err
	}

	return NewLogger(
		config.Intercept,
		encoder,
		st,
		logger,
	)
}

// NewLogger creates a new audit logging pipeline with the provided elements.
func NewLogger(
	intercept InterceptConfig,
	encoder codec.Encoder,
	storage storage.WritableStorage,
	logger log.Logger,
) (Logger, error) {
	return &loggerImplementation{
		intercept: intercept,
		encoder:   encoder,
		storage:   storage,
		logger:    logger,
	}, nil
}

// NewEncoder creates a new audit log encoder of the specified format.
func NewEncoder(encoder Format, logger log.Logger) (codec.Encoder, error) {
	switch encoder {
	case FormatNone:
		return noneCodec.NewEncoder(), nil
	case FormatAsciinema:
		return asciinema.NewEncoder(logger), nil
	case FormatBinary:
		return binary.NewEncoder(), nil
	default:
		return nil, fmt.Errorf("invalid audit log encoder: %s", encoder)
	}
}

// NewStorage creates a new audit log storage of the specified type and with the specified configuration.
func NewStorage(config Config, logger log.Logger) (storage.WritableStorage, error) {
	switch config.Storage {
	case StorageNone:
		return noneStorage.NewStorage(), nil
	case StorageFile:
		return file.NewStorage(config.File, logger)
	case StorageS3:
		return s3.NewStorage(config.S3, logger)
	default:
		return nil, fmt.Errorf("invalid audit log storage: %s", config.Storage)
	}
}
