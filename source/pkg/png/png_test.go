package png

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/Mihai22125/SteganoGo/pkg/fileutils"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type testCase struct {
	testName      string
	fileName      string
	expectedData  string
	expectedError error
}

type testGroup struct {
	groupName   string
	testDir     string
	expectedDir string
	testFiles   []string
	testCases   []testCase
}

type MySuite struct {
	testPath    string
	expectedDir string
	shouldPass  []testGroup
	shouldFail  []testGroup
	testPNG     StructPNG
}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpSuite(c *C) {
	s.testPath = "../../../test_files/png_test_files/test_files"
	s.expectedDir = "../../../test_files/png_test_files/expected/expected_parsed_json"

	s.shouldPass = []testGroup{
		{
			groupName:   "basic",
			testDir:     path.Join(s.testPath, "01.basic_formats"),
			expectedDir: path.Join(s.expectedDir, "01.basic_formats"),
		},
	}

	for i, sp := range s.shouldPass {

		files, err := ioutil.ReadDir(sp.expectedDir)
		c.Assert(err, IsNil)

		for j, file := range files {
			s.shouldPass[i].testFiles = append(s.shouldPass[i].testFiles, fileutils.FileNameWithoutExtension(file.Name()))

			expected, err := ioutil.ReadFile(path.Join(sp.expectedDir, file.Name()))
			expected = bytes.Replace(expected, []byte("\x0d"), []byte{}, -1)
			c.Assert(err, IsNil)
			newTestCase := testCase{
				testName:      sp.groupName + fmt.Sprintf("%d", j),
				fileName:      fileutils.FileNameWithoutExtension(file.Name()),
				expectedData:  string(expected),
				expectedError: nil,
			}

			s.shouldPass[i].testCases = append(s.shouldPass[i].testCases, newTestCase)

		}
	}

	s.shouldFail = []testGroup{
		{
			groupName: "corruptedFile",
			testDir:   path.Join(s.testPath, "13.corrupted_files2"),
		},
	}

	//s.shouldFail[0].testFiles = []string{"cor201.png"}
	s.shouldFail[0].testCases = []testCase{
		{
			testName:      "PNG file without header",
			fileName:      "cor201",
			expectedData:  "",
			expectedError: ErrBadPNG,
		},
		{
			testName:      "Invalid PNG header",
			fileName:      "cor202",
			expectedData:  "",
			expectedError: ErrNotPNG,
		},
		{
			testName:      "Invalid Chunk size",
			fileName:      "cor203",
			expectedData:  "",
			expectedError: ErrPNGChunks,
		},
	}

	s.testPNG = StructPNG{
		header: Header{
			signature: []uint8{0x89, 0x50, 0x4e, 0x47, 0xd, 0xa, 0x1a, 0xa},
		},
		chunks: []Chunk{
			{
				size:      0xd,
				chunkType: "IHDR",
				data:      []uint8{0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x20, 0x1, 0x0, 0x0, 0x0, 0x0},
				crc:       0x5b014759,
			},
			{
				size:      0x5b,
				chunkType: "IDAT",
				data:      []uint8{0x78, 0x9c, 0x2d, 0xcc, 0xb1, 0x9, 0x3, 0x30, 0xc, 0x5, 0xd1, 0xeb, 0xd2, 0x4, 0xb2, 0x4a, 0x20, 0xb, 0x7a, 0x34, 0x6f},
				crc:       0xd02f14c9,
			},
			{
				size:      0x5b,
				chunkType: "IDAT",
				data:      []uint8{0x90, 0x15, 0x3c, 0x82, 0xc1, 0x8d, 0xa, 0x61, 0x45, 0x7, 0x51, 0xf1, 0xe0, 0x8a, 0x2f, 0xaa, 0xea, 0xd2, 0xa4, 0x84},
				crc:       0xd02f14c9,
			},
			{
				size:      0x0,
				chunkType: "IEND",
				data:      []uint8{},
				crc:       0xae426082},
		},
	}

}

func (s *MySuite) TestparsePNG(c *C) {

	for _, sp := range s.shouldPass {
		for _, testCase := range sp.testCases {

			testFile, err := os.Open(path.Join(sp.testDir, testCase.fileName) + ".png")
			c.Assert(err, IsNil)
			defer testFile.Close()

			result, err := ParsePNG(testFile)
			c.Assert(err, Equals, testCase.expectedError)

			json, err := json.MarshalIndent(result, "", "  ")
			c.Assert(err, IsNil)

			c.Assert(string(json), Equals, testCase.expectedData)

		}
	}

	for _, sp := range s.shouldFail {
		for _, testCase := range sp.testCases {
			testFile, err := os.Open(path.Join(sp.testDir, testCase.fileName) + ".png")
			c.Assert(err, IsNil)
			defer testFile.Close()

			result, err := ParsePNG(testFile)
			c.Assert(err, Equals, testCase.expectedError)
			c.Assert(result, DeepEquals, EmptyPNG)

		}
	}
}

func (s *MySuite) TestParsePNG(c *C) {

	// should pass
	c.Assert(s.shouldPass, NotNil)
	c.Assert(s.shouldPass[0].testCases, NotNil)

	testFile, err := os.Open(path.Join(s.shouldPass[0].testDir, s.shouldPass[0].testCases[0].fileName) + ".png")
	c.Assert(err, IsNil)
	defer testFile.Close()

	result, err := ParsePNG(testFile)
	c.Assert(err, IsNil)

	json, err := json.MarshalIndent(result, "", "  ")
	c.Assert(err, IsNil)
	c.Assert(string(json), Equals, s.shouldPass[0].testCases[0].expectedData)

	// should fail - nil pointer refference for *os.File Argument

	result, err = ParsePNG(nil)
	c.Assert(err, Equals, os.ErrInvalid)
	c.Assert(result, DeepEquals, EmptyPNG)
}

func (s *MySuite) TestreadSingleChunk(c *C) {

	tt := []struct {
		name           string
		rawData        []byte
		expectedResult Chunk
		expectedErr    error
	}{
		{
			name:    "ShouldPass1",
			rawData: []byte{0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x20, 0x01, 0x00, 0x00, 0x00, 0x00, 0x5B, 0x01, 0x47, 0x59},
			expectedResult: Chunk{
				size:      13,
				chunkType: TypeIHDR,
				data:      []uint8{0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x20, 0x1, 0x0, 0x0, 0x0, 0x0},
				crc:       1526810457,
			},
			expectedErr: nil,
		},
		{
			name:           "ShouldFail1",
			rawData:        []byte{0x00, 0x00, 0x00},
			expectedResult: Chunk{},
			expectedErr:    io.ErrUnexpectedEOF,
		},
		{
			name:           "ShouldFail2",
			rawData:        []byte{0x00, 0x00, 0x00, 0x0D, 0x49, 0x48},
			expectedResult: Chunk{},
			expectedErr:    io.ErrUnexpectedEOF,
		},
		{
			name:           "ShouldFail3",
			rawData:        []byte{0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, 0x00, 0x00},
			expectedResult: Chunk{},
			expectedErr:    io.ErrUnexpectedEOF,
		},
		{
			name:           "ShouldFail4",
			rawData:        []byte{0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x20, 0x01, 0x00, 0x00, 0x00, 0x00, 0x5B, 0x01},
			expectedResult: Chunk{},
			expectedErr:    io.ErrUnexpectedEOF,
		},
	}

	for _, tc := range tt {
		myReader := bytes.NewReader(tc.rawData)
		result, err := readSingleChunck(myReader)
		c.Assert(result, DeepEquals, tc.expectedResult)
		c.Assert(err, Equals, tc.expectedErr)
	}

}

func (s *MySuite) TestreadChunks(c *C) {

	tt := []struct {
		name           string
		rawData        []byte
		expectedResult []Chunk
		expectedErr    error
	}{
		{
			name: "ShouldPass",
			rawData: []byte{
				0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x20,
				0x01, 0x00, 0x00, 0x00, 0x00, 0x5B, 0x01, 0x47, 0x59, 0x00, 0x00, 0x00, 0x04, 0x67, 0x41, 0x4D,
				0x41, 0x00, 0x01, 0x86, 0xA0, 0x31, 0xE8, 0x96, 0x5F, 0x00, 0x00, 0x00, 0x5B, 0x49, 0x44, 0x41,
				0x54, 0x78, 0x9C, 0x2D, 0xCC, 0xB1, 0x09, 0x03, 0x30, 0x0C, 0x05, 0xD1, 0xEB, 0xD2, 0x04, 0xB2,
				0x4A, 0x20, 0x0B, 0x7A, 0x34, 0x6F, 0x90, 0x15, 0x3C, 0x82, 0xC1, 0x8D, 0x0A, 0x61, 0x45, 0x07,
				0x51, 0xF1, 0xE0, 0x8A, 0x2F, 0xAA, 0xEA, 0xD2, 0xA4, 0x84, 0x6C, 0xCE, 0xA9, 0x25, 0x53, 0x06,
				0xE7, 0x53, 0x34, 0x57, 0x12, 0xE2, 0x11, 0xB2, 0x21, 0xBF, 0x4B, 0x26, 0x3D, 0x1B, 0x42, 0x73,
				0x25, 0x25, 0x5E, 0x8B, 0xDA, 0xB2, 0x9E, 0x6F, 0x6A, 0xCA, 0x30, 0x69, 0x2E, 0x9D, 0x29, 0x61,
				0x6E, 0xE9, 0x6F, 0x30, 0x65, 0xF0, 0xBF, 0x1F, 0x10, 0x87, 0x49, 0x2F, 0xD0, 0x2F, 0x14, 0xC9,
				0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82},
			expectedResult: []Chunk{
				{
					size:      0xd,
					chunkType: "IHDR",
					data:      []uint8{0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x20, 0x1, 0x0, 0x0, 0x0, 0x0},
					crc:       0x5b014759,
				},
				{
					size:      0x4,
					chunkType: "gAMA",
					data:      []uint8{0x0, 0x1, 0x86, 0xa0},
					crc:       0x31e8965f,
				},
				{size: 0x5b,
					chunkType: "IDAT",
					data:      []uint8{0x78, 0x9c, 0x2d, 0xcc, 0xb1, 0x9, 0x3, 0x30, 0xc, 0x5, 0xd1, 0xeb, 0xd2, 0x4, 0xb2, 0x4a, 0x20, 0xb, 0x7a, 0x34, 0x6f, 0x90, 0x15, 0x3c, 0x82, 0xc1, 0x8d, 0xa, 0x61, 0x45, 0x7, 0x51, 0xf1, 0xe0, 0x8a, 0x2f, 0xaa, 0xea, 0xd2, 0xa4, 0x84, 0x6c, 0xce, 0xa9, 0x25, 0x53, 0x6, 0xe7, 0x53, 0x34, 0x57, 0x12, 0xe2, 0x11, 0xb2, 0x21, 0xbf, 0x4b, 0x26, 0x3d, 0x1b, 0x42, 0x73, 0x25, 0x25, 0x5e, 0x8b, 0xda, 0xb2, 0x9e, 0x6f, 0x6a, 0xca, 0x30, 0x69, 0x2e, 0x9d, 0x29, 0x61, 0x6e, 0xe9, 0x6f, 0x30, 0x65, 0xf0, 0xbf, 0x1f, 0x10, 0x87, 0x49, 0x2f},
					crc:       0xd02f14c9,
				},
				{
					size:      0x0,
					chunkType: "IEND",
					data:      []uint8{},
					crc:       0xae426082,
				},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range tt {
		myReader := bytes.NewReader(tc.rawData)
		result, err := readChunks(myReader)
		c.Assert(result, DeepEquals, tc.expectedResult)
		c.Assert(err, Equals, tc.expectedErr)
	}
}

func (s *MySuite) TestIHDRChunk(c *C) {

	myPNG := s.testPNG
	// shouldPass

	ihdrChunk, err := myPNG.IHDRChunk()
	c.Assert(err, IsNil)
	c.Assert(ihdrChunk, DeepEquals, myPNG.chunks[0])

	// delete ihdr chunk
	myPNG.chunks = myPNG.chunks[1:]

	// shouldFail
	ihdrChunk, err = myPNG.IHDRChunk()
	c.Assert(err, Equals, ErrIHDRMissing)
	c.Assert(ihdrChunk, DeepEquals, Chunk{})
}

func (s *MySuite) TestIDATChunks(c *C) {

	myPNG := s.testPNG
	// shouldPass

	idatChunks, err := myPNG.IDATChunks()
	c.Assert(err, IsNil)
	c.Assert(idatChunks, DeepEquals, []Chunk{myPNG.chunks[1], myPNG.chunks[2]})

	// delete IDAT chunks
	myPNG.chunks = append(myPNG.chunks[:1], myPNG.chunks[1+2:]...)
	idatChunks, err = myPNG.IDATChunks()
	c.Assert(err, Equals, ErrIDATMissing)
	c.Assert(idatChunks, IsNil)
}

func (s *MySuite) TestIDATData(c *C) {

	myPNG := s.testPNG
	// shouldPass

	expectedIDATdata := []byte{0x90, 0x15, 0x3c, 0x82, 0xc1, 0x8d, 0xa, 0x61, 0x45, 0x7, 0x51, 0xf1, 0xe0, 0x8a, 0x2f, 0xaa, 0xea, 0xd2, 0xa4, 0x84}
	idatData, err := myPNG.IDATdata()
	c.Assert(err, IsNil)
	c.Assert(idatData, DeepEquals, expectedIDATdata)

	// delete IDAT chunks
	myPNG.chunks = append(myPNG.chunks[:1], myPNG.chunks[1+2:]...)

	idatData, err = myPNG.IDATdata()
	c.Assert(err, Equals, ErrIDATMissing)
	c.Assert(idatData, IsNil)
}

// TestPNGHeader test for StructPNG.Header() method
func (s *MySuite) TestPNGHeader(c *C) {
	myPNG := s.testPNG

	pngHeader := myPNG.Header()

	c.Assert(pngHeader, DeepEquals, Header{[]byte{137, 80, 78, 71, 13, 10, 26, 10}})
}

// test Chunk getters

func (s *MySuite) TestChunkSize(c *C) {
	myChunk := Chunk{
		size:      13,
		chunkType: TypeIHDR,
		data:      []uint8{0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x20, 0x1, 0x0, 0x0, 0x0, 0x0},
		crc:       1526810457,
	}

	size := myChunk.Size()
	c.Assert(size, Equals, myChunk.size)
}

func (s *MySuite) TestChunkType(c *C) {
	myChunk := Chunk{
		size:      13,
		chunkType: TypeIHDR,
		data:      []uint8{0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x20, 0x1, 0x0, 0x0, 0x0, 0x0},
		crc:       1526810457,
	}

	Ctype := myChunk.ChunkType()
	c.Assert(Ctype, Equals, myChunk.chunkType)
}

func (s *MySuite) TestChunkData(c *C) {
	myChunk := Chunk{
		size:      13,
		chunkType: TypeIHDR,
		data:      []uint8{0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x20, 0x1, 0x0, 0x0, 0x0, 0x0},
		crc:       1526810457,
	}

	Cdata := myChunk.Data()
	c.Assert(Cdata, DeepEquals, myChunk.data)
}

func (s *MySuite) TestChunkCRC(c *C) {
	myChunk := Chunk{
		size:      13,
		chunkType: TypeIHDR,
		data:      []uint8{0x0, 0x0, 0x0, 0x20, 0x0, 0x0, 0x0, 0x20, 0x1, 0x0, 0x0, 0x0, 0x0},
		crc:       1526810457,
	}

	Ccrc := myChunk.CRC()
	c.Assert(Ccrc, Equals, myChunk.crc)
}
