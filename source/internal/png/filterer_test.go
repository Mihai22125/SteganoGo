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
	s.filterer = newFilterer(1, 32, 32, 8)

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

func (s *MySuiteFilterer) TestAbs8(c *C) {
	c.Assert(abs8(10), Equals, 10)
	c.Assert(abs8(138), Equals, 118)
}

func (s *MySuiteFilterer) TestFilterData(c *C) {
	filterer := newFilterer(1, 8, 4, 1)
	rawData := []byte{255, 0, 8, 255, 8, 15, 255, 16, 23, 255, 24, 31, 255, 32, 39, 255, 41, 47, 255, 49, 55, 255, 57, 63, 255, 65, 71, 255, 74, 79, 255, 82}
	filteredData := filterer.FilterData(rawData)

	expected := []byte{0x0, 0xff, 0x0, 0x8, 0xff, 0x8, 0xf, 0xff, 0x10, 0x0, 0x17, 0xff, 0x18, 0x1f, 0xff, 0x20, 0x27, 0xff, 0x0, 0x29, 0x2f, 0xff, 0x31, 0x37, 0xff, 0x39, 0x3f, 0x0, 0xff, 0x41, 0x47, 0xff, 0x4a, 0x4f, 0xff, 0x52}
	c.Assert(filteredData, DeepEquals, expected)
}
