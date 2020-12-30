package s3

import (
	"fmt"
	"os"
)

// Config S3 storage configuration
type Config struct {
	Local           string   `json:"local" yaml:"local" default:"/var/lib/audit"`
	AccessKey       string   `json:"accessKey" yaml:"accessKey"`
	SecretKey       string   `json:"secretKey" yaml:"secretKey"`
	Bucket          string   `json:"bucket" yaml:"bucket"`
	Region          string   `json:"region" yaml:"region"`
	Endpoint        string   `json:"endpoint" yaml:"endpoint"`
	CACert          string   `json:"cacert" yaml:"cacert"`
	ACL             string   `json:"acl" yaml:"acl"`
	PathStyleAccess bool     `json:"pathStyleAccess" yaml:"pathStyleAccess"`
	UploadPartSize  uint     `json:"uploadPartSize" yaml:"uploadPartSize" default:"5242880"`
	ParallelUploads uint     `json:"parallelUploads" yaml:"parallelUploads" default:"20"`
	Metadata        Metadata `json:"metadata" yaml:"metadata"`
}

// Validate validates the
func (config Config) Validate() error {
	if config.Local == "" {
		return fmt.Errorf("empty local storage directory provided")
	}
	stat, err := os.Stat(config.Local)
	if err != nil {
		return fmt.Errorf("invalid local directory: %s (%w)", config.Local, err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("invalid local directory: %s (not a directory)", config.Local)
	}
	if config.AccessKey == "" {
		return fmt.Errorf("no access key provided")
	}
	if config.SecretKey == "" {
		return fmt.Errorf("no secret key provided")
	}
	if config.Bucket == "" {
		return fmt.Errorf("no bucket name provided")
	}
	if config.UploadPartSize < 5242880 {
		return fmt.Errorf("upload part size too low %d (minimum 5 MB)", config.UploadPartSize)
	}
	if config.ParallelUploads < 1 {
		return fmt.Errorf("parallel uploads invalid: %d (must be positive)", config.ParallelUploads)
	}
	return nil
}

// Metadata Metadata configuration for the S3 storage
type Metadata struct {
	IP       bool `json:"ip" yaml:"ip"`
	Username bool `json:"username" yaml:"username"`
}
