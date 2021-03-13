package fileutils

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

// PreProcessFile consume a file handle of type *os.File and return a type of *bytes.Reader
func PreProcessFile(dat *os.File) (*bytes.Reader, error) {
	stats, err := dat.Stat()
	if err != nil {
		return nil, err
	}

	var size = stats.Size()
	b := make([]byte, size)

	buf := bufio.NewReader(dat)
	_, err = buf.Read(b)
	if err != nil {
		return nil, err
	}

	bReader := bytes.NewReader(b)

	return bReader, nil
}

// FileNameWithoutWxtension returns filename without extension
func FileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
