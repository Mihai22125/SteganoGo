package fileutils

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.

func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
	Dir  string // temporary directory for testing
	f    string // Temporary file for testing
	text []byte
}

var _ = Suite(&MySuite{})

// Setupsuite prepares temporary files for testing
func (s *MySuite) SetUpSuite(c *C) {
	text := []byte("hello")
	dir := c.MkDir() // The directory created by c.MkDir() will be automatically destroyed after the Suite ends.

	tmpfile, err := ioutil.TempFile(dir, "")
	if err != nil {
		c.Errorf("Fail to create test file: %v\n", tmpfile.Name(), err)
	}
	defer tmpfile.Close()

	_, err = tmpfile.Write(text)
	if err != nil {
		c.Errorf("Fail to prepare test file.%v\n", tmpfile.Name(), err)
	}
	s.text = text
	s.Dir = dir
	s.f = filepath.Base(tmpfile.Name())
}

func (s *MySuite) TestparsePNG(c *C) {
	// should pass
	file, err := os.Open(path.Join(s.Dir, s.f))
	c.Assert(err, IsNil)
	defer file.Close()

	buf, err := PreProcessFile(file)
	c.Assert(err, IsNil)

	data := make([]byte, len(s.text))
	n, err := buf.Read(data)
	c.Assert(err, IsNil)
	c.Assert(n, Equals, len(s.text))
	c.Assert(data, DeepEquals, s.text)

	// should fail, nil argument
	buf, err = PreProcessFile(nil)
	c.Assert(err, Equals, os.ErrInvalid)
	c.Assert(buf, IsNil)

}

func (s *MySuite) TestFileNameWithoutExtension(c *C) {
	tt := []struct {
		name           string
		input          string
		expectedOutput string
	}{
		{name: "test1", input: "file1.png", expectedOutput: "file1"},
		{name: "test2", input: "file2", expectedOutput: "file2"},
		{name: "test3", input: "file3.tar.gz", expectedOutput: "file3.tar"},
		{name: "test4", input: "", expectedOutput: ""},
	}

	for _, tc := range tt {
		result := FileNameWithoutExtension(tc.input)
		c.Assert(result, Equals, tc.expectedOutput)
	}

}
