package file

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/containerssh/auditlog/storage"
)

type Storage struct {
	directory string
}

func (s *Storage) OpenReader(name string) (io.Reader, error) {
	return os.Open(path.Join(s.directory, name))
}

func (s *Storage) List() (<-chan storage.Entry, <-chan error) {
	result := make(chan storage.Entry)
	errorChannel := make(chan error)
	go func() {
		if err := filepath.Walk(s.directory, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() && info.Size() > 0 && !strings.Contains(info.Name(), ".") {
				result <- storage.Entry{
					Name:     info.Name(),
					Metadata: map[string]string{},
				}
			}
			return err
		}); err != nil {
			errorChannel <- err
		}
		close(result)
		close(errorChannel)
	}()
	return result, errorChannel
}

func (s *Storage) OpenWriter(name string) (storage.Writer, error) {
	file, err := os.Create(path.Join(s.directory, name))
	if err != nil {
		return nil, err
	}
	return &writer{
		file: file,
	}, nil

}

type writer struct {
	file *os.File
}

func (w *writer) Write(p []byte) (n int, err error) {
	return w.file.Write(p)
}

func (w *writer) Close() error {
	return w.file.Close()
}

func (w *writer) SetMetadata(_ int64, _ string, _ *string) {

}
