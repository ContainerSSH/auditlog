package auditlog

import (
	"github.com/containerssh/auditlog/storage/file"
	"github.com/containerssh/auditlog/storage/s3"
)

// Format describes the audit log format in use.
type Format string

const (
	// FormatNone signals that no audit logging should take place.
	FormatNone Format = "none"
	// FormatBinary signals that audit logging should take place in CBOR+GZIP format
	//              (see https://containerssh.github.io/advanced/audit/format/ )
	FormatBinary Format = "binary"
	// FormatAsciinema signals that audit logging should take place in Asciicast v2 format
	//                 (see https://github.com/asciinema/asciinema/blob/develop/doc/asciicast-v2.md )
	FormatAsciinema Format = "asciinema"
)

// Storage describes the storage backend to use.
type Storage string

const (
	// StorageNone signals that no storage should be used.
	StorageNone Storage = "none"
	// StorageFile signals that audit logs should be stored in a local directory.
	StorageFile Storage = "file"
	// StorageS3 signals that audit logs should be stored in an S3-compatible object storage.
	StorageS3 Storage = "s3"
)

// Config is the configuration structure for audit logging.
type Config struct {
	// Format audit format
	Format Format `json:"format" yaml:"format" default:"none"`
	// Storage audit storage type
	Storage Storage `json:"storage" yaml:"storage" default:"none"`
	// File audit logger configuration
	File file.Config `json:"file" yaml:"file"`
	// S3 configuration
	S3 s3.Config `json:"s3" yaml:"s3"`
	// Intercept configures what should be intercepted
	Intercept InterceptConfig `json:"intercept" yaml:"intercept"`
}

// InterceptConfig configures what should be intercepted by the auditing facility.
type InterceptConfig struct {
	// Stdin signals that the standard input from the user should be captured.
	Stdin bool `json:"stdin" yaml:"stdin" default:"false"`
	// Stdout signals that the standard output to the user should be captured.
	Stdout bool `json:"stdout" yaml:"stdout" default:"false"`
	// Stderr signals that the standard error to the user should be captured.
	Stderr bool `json:"stderr" yaml:"stderr" default:"false"`
	// Passwords signals that passwords during authentication should be captured.
	Passwords bool `json:"passwords" yaml:"passwords" default:"false"`
}
