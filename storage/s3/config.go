package s3

type Config struct {
	Local           string   `json:"local" yaml:"local" default:"/var/lib/audit"`
	AccessKey       string   `json:"accessKey" yaml:"accessKey"`
	SecretKey       string   `json:"secretKey" yaml:"secretKey"`
	Bucket          string   `json:"bucket" yaml:"bucket"`
	Region          string   `json:"region" yaml:"region"`
	Endpoint        string   `json:"endpoint" yaml:"endpoint"`
	CaCert          string   `json:"cacert" yaml:"cacert"`
	ACL             string   `json:"acl" yaml:"acl"`
	UploadPartSize  uint     `json:"uploadPartSize" yaml:"uploadPartSize"`
	ParallelUploads uint     `json:"parallelUploads" yaml:"parallelUploads"`
	Metadata        Metadata `json:"metadata" yaml:"metadata"`
}

type Metadata struct {
	IP       bool `json:"ip" yaml:"ip"`
	Username bool `json:"username" yaml:"username"`
}
