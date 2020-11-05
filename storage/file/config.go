package file

// Config A configuration for the file storage.
type Config struct {
	Directory string `json:"directory" yaml:"directory" default:"/var/log/audit"`
}
