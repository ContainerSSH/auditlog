package file

import (
	"fmt"
	"os"
)

// Config is the configuration for the file storage.
type Config struct {
	Directory string `json:"directory" yaml:"directory" default:"/var/log/audit"`
}

func (c *Config) Validate() error {
	stat, err := os.Stat(c.Directory)
	if err != nil {
		return fmt.Errorf("invalid audit log directory: %s (%w)", c.Directory, err)
	}
	if !stat.IsDir() {
		return fmt.Errorf("invalid audit log directory: %s (not a directory)", c.Directory)
	}
	return nil
}
