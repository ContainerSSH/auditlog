package auditlog

import (
	"github.com/containerssh/auditlog/storage/file"
	"github.com/containerssh/auditlog/storage/s3"
)

// swagger:enum Format
type Format string

const (
	FormatNone      Format = "none"
	FormatBinary    Format = "binary"
	FormatAsciinema Format = "asciinema"
)

// swagger:enum Storage
type Storage string

const (
	StorageNone Storage = "none"
	StorageFile Storage = "file"
	StorageS3   Storage = "s3"
)

type Config struct {
	// Audit format
	Format Format `json:"format" yaml:"format" default:"none"`
	// Audit storage type
	Storage Storage `json:"storage" yaml:"storage" default:"none"`
	// File audit logger configuration
	File file.Config `json:"file" yaml:"file"`
	// S3 configuration
	S3 s3.Config `json:"s3" yaml:"s3"`
	// What to intercept during the connection
	Intercept InterceptConfig `json:"intercept" yaml:"intercept"`
}

type InterceptConfig struct {
	Stdin     bool `json:"stdin" yaml:"stdin" default:"false"`
	Stdout    bool `json:"stdout" yaml:"stdout" default:"false"`
	Stderr    bool `json:"stderr" yaml:"stderr" default:"false"`
	Passwords bool `json:"passwords" yaml:"passwords" default:"false"`
}
