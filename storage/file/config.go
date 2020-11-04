package file

type Config struct {
	Directory string `json:"directory" yaml:"directory" default:"/var/log/audit"`
}
