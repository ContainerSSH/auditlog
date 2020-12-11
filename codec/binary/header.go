package binary

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const fileFormat string = "ContainerSSH-Auditlog"
const fileFormatLength = 32
const currentVersion = uint64(1)

var fileFormatBytes []byte

type header struct {
	fileFormat []byte
	version    uint64
}

func (h header) getBytes() []byte {
	versionBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(versionBytes, h.version)
	result := make([]byte, fileFormatLength+8)
	for i := 0; i < fileFormatLength; i++ {
		result[i] = fileFormatBytes[i]
	}
	for i := 0; i < 8; i++ {
		result[i+fileFormatLength] = versionBytes[i]
	}
	return result
}

func newHeader(version uint64) header {
	return header{
		fileFormat: fileFormatBytes,
		version:    version,
	}
}

func readHeader(reader io.Reader, maxVersion uint64) error {
	headerBytes := make([]byte, fileFormatLength+8)

	_, err := io.ReadAtLeast(reader, headerBytes, fileFormatLength+8)
	if err != nil {
		return err
	}
	if !bytes.Equal(headerBytes[:fileFormatLength], fileFormatBytes) {
		return fmt.Errorf("invalid file format header: %v", headerBytes[:fileFormatLength])
	}
	version := binary.LittleEndian.Uint64(headerBytes[fileFormatLength:])
	if version > maxVersion {
		return fmt.Errorf("file format version is higher than supported: %d", version)
	}
	return nil
}

func init() {
	fileFormatBytes = make([]byte, fileFormatLength)
	for i := 0; i < len(fileFormat); i++ {
		fileFormatBytes[i] = fileFormat[i]
	}
	for i := len(fileFormat); i < fileFormatLength; i++ {
		fileFormatBytes[i] = "\000"[0]
	}
}
