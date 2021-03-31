package pngint

import (
	"io/ioutil"
	"path/filepath"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
//func Test(t *testing.T) { TestingT(t) }

type testCase struct {
	name           string
	inputFileName  string
	outputFileName string
	input          []byte
	expectedOutput []byte
	expectedError  error
}

type MySuiteFilterer struct {
	filterer   *Filterer
	shouldPass []testCase
	shouldFail []testCase
	inputDir   string
	outputDir  string
}

var _ = Suite(&MySuiteFilterer{})

func (s *MySuiteFilterer) SetUpSuite(c *C) {
	s.filterer = newFilterer(1, 32, 32)

	s.inputDir = "../../../test_files/png_test_files/png_internal/filtering/input"
	s.outputDir = "../../../test_files/png_test_files/png_internal/filtering/expected"

	s.shouldPass = []testCase{
		{
			name:           "no filtering",
			inputFileName:  "noFiltering.in",
			outputFileName: "noFiltering.in",
			expectedError:  nil,
		},
		{
			name:           "sub filtering",
			inputFileName:  "subFiltering.in",
			outputFileName: "subFiltering.in",
			expectedError:  nil,
		},
		{
			name:           "up filtering",
			inputFileName:  "upFiltering.in",
			outputFileName: "upFiltering.in",
			expectedError:  nil,
		},
		{
			name:           "averagefiltering",
			inputFileName:  "averageFiltering.in",
			outputFileName: "averageFiltering.in",
			expectedError:  nil,
		},
		{
			name:           "paeth filtering",
			inputFileName:  "paethFiltering.in",
			outputFileName: "paethFiltering.in",
			expectedError:  nil,
		},
	}

	for i, tc := range s.shouldPass {
		dat, err := ioutil.ReadFile(filepath.Join(s.inputDir, tc.inputFileName))
		c.Assert(err, IsNil)
		s.shouldPass[i].input = dat
		dat, err = ioutil.ReadFile(filepath.Join(s.outputDir, tc.outputFileName))
		c.Assert(err, IsNil)
		s.shouldPass[i].expectedOutput = dat
	}

	s.shouldFail = []testCase{
		{
			name:           "empty imput slice",
			input:          []byte{},
			expectedOutput: []byte{},
			expectedError:  ErrInvalidInput,
		},
		{
			name:           "nil imput slice",
			input:          nil,
			expectedOutput: []byte{},
			expectedError:  ErrInvalidInput,
		},
		{
			name:           "unknown filtering type",
			input:          []byte{0x05, 0x23, 0x43, 0x00, 0x00},
			expectedOutput: []byte{},
			expectedError:  ErrUnknownFilterType,
		},
	}

}

func (s *MySuiteFilterer) Testreconstruct(c *C) {

	for _, tc := range s.shouldPass {
		err := s.filterer.reconstruct(tc.input)
		c.Assert(err, Equals, tc.expectedError)
		c.Assert(s.filterer.recon, DeepEquals, tc.expectedOutput)
	}

	for _, tc := range s.shouldFail {
		err := s.filterer.reconstruct(tc.input)
		c.Assert(err, Equals, tc.expectedError)
		c.Assert(s.filterer.recon, DeepEquals, tc.expectedOutput)
	}

}
