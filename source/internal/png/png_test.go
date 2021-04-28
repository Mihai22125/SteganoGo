package pngint

import (
	"encoding/json"
	"io/ioutil"
	"path"

	"github.com/Mihai22125/SteganoGo/pkg/png"
	. "gopkg.in/check.v1"
)

type MySuitePNG struct {
	myPNG *pngImage
}

var _ = Suite(&MySuitePNG{})

func (s *MySuitePNG) SetUpSuite(c *C) {
	s.myPNG = new(pngImage)
}

func (s *MySuitePNG) TestProcessIHDR(c *C) {

	testTable := []struct {
		name         string
		IHDRchunk    []byte
		expectedMeta imageMetadata
		expectedErr  error
	}{
		{
			name:         "shouldPass1 - real smaple",
			IHDRchunk:    []byte{0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x20, 0x08, 0x00, 0x00, 0x00, 0x00},
			expectedMeta: imageMetadata{width: 32, height: 32, bitDepth: 8, colorType: 0, compressionMethod: 0, filterMethod: 0, interlaceMethod: 0},
			expectedErr:  nil,
		},

		{
			name:         "shouldFail - invalid IHDR length",
			IHDRchunk:    []byte{0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x20, 0x08, 0x00, 0x00, 0x00},
			expectedMeta: imageMetadata{},
			expectedErr:  png.ErrInvalidIHDR,
		},
	}

	for _, testCase := range testTable {
		err := s.myPNG.processIHDR(testCase.IHDRchunk)
		c.Assert(err, Equals, testCase.expectedErr)
		c.Assert(s.myPNG.meta, DeepEquals, testCase.expectedMeta)
	}

}

func (s *MySuitePNG) TestExtractMetadata(c *C) {

	testDir := "../../../test_files/png_test_files/png_internal/extractMetadata/input"
	_ = testDir
	testTable := []struct {
		testName     string
		test_file    string
		png          png.StructPNG
		expectedMeta imageMetadata
		expectedErr  error
	}{
		{
			testName:     "ShouldPass1",
			test_file:    "shouldPass1.txt",
			expectedMeta: imageMetadata{width: 32, height: 32, bitDepth: 1, colorType: 0, compressionMethod: 0, filterMethod: 0, interlaceMethod: 0},
			expectedErr:  nil,
		},
		{
			testName:     "ShouldFail1 - no IHDR chunk",
			test_file:    "shouldFail1_no_IHDR.txt",
			expectedMeta: imageMetadata{},
			expectedErr:  png.ErrIHDRMissing,
		},

		{
			testName:     "ShouldFail1 - invalid IHDR chunk",
			test_file:    "shouldFail2_Invalid_IHDR.txt",
			expectedMeta: imageMetadata{},
			expectedErr:  png.ErrInvalidIHDR,
		},
	}
	// setup test cases
	for i, testCase := range testTable {
		testData, err := ioutil.ReadFile(path.Join(testDir, testCase.test_file))
		c.Assert(err, IsNil)

		png := png.StructPNG{}

		err = json.Unmarshal(testData, &png)
		c.Assert(err, IsNil)

		testTable[i].png = png
	}

	for _, testCase := range testTable {
		err := s.myPNG.extractMetadata(testCase.png)
		c.Assert(err, Equals, testCase.expectedErr)
		c.Assert(s.myPNG.meta, DeepEquals, testCase.expectedMeta)
	}
}

func (s *MySuitePNG) TestbytesPerPixel(c *C) {

	testTable := []struct {
		testName       string
		png            pngImage
		expectedResult uint8
	}{
		{
			testName: "ShouldPass1 -  less than 1 byte per sample, each pixel contains one sample",
			png: pngImage{
				meta: imageMetadata{width: 32, height: 32, bitDepth: 1, colorType: 0, compressionMethod: 0, filterMethod: 0, interlaceMethod: 0},
			},
			expectedResult: 1,
		},
		{
			testName: "ShouldPass2 - one byte per sample, each pixel contains one sample",
			png: pngImage{
				meta: imageMetadata{width: 32, height: 32, bitDepth: 8, colorType: 0, compressionMethod: 0, filterMethod: 0, interlaceMethod: 0},
			},
			expectedResult: 1,
		},
		{
			testName: "ShouldPass3 - more than one byte per sample, each pixel contains one sample",
			png: pngImage{
				meta: imageMetadata{width: 32, height: 32, bitDepth: 16, colorType: 0, compressionMethod: 0, filterMethod: 0, interlaceMethod: 0},
			},
			expectedResult: 2,
		},
		{
			testName: "ShouldPass4 - one byte per sample, each pixel contains more than one sample",
			png: pngImage{
				meta: imageMetadata{width: 32, height: 32, bitDepth: 8, colorType: 2, compressionMethod: 0, filterMethod: 0, interlaceMethod: 0},
			},
			expectedResult: 3,
		},
		{
			testName: "ShouldPass5 - more than one byte per sample, each pixel contains more than one sample",
			png: pngImage{
				meta: imageMetadata{width: 32, height: 32, bitDepth: 16, colorType: 2, compressionMethod: 0, filterMethod: 0, interlaceMethod: 0},
			},
			expectedResult: 6,
		},
	}

	for _, testCase := range testTable {
		bytesPerPixel := testCase.png.bytesPerPixel()
		c.Assert(bytesPerPixel, Equals, testCase.expectedResult)
	}
}

func (s *MySuitePNG) TestSamplesPerPixel(c *C) {

	testTable := []struct {
		testName       string
		png            pngImage
		expectedResult uint8
	}{
		{
			testName: "color type 0 - Each pixel is a grayscale sample",
			png: pngImage{
				meta: imageMetadata{width: 32, height: 32, bitDepth: 1, colorType: 0, compressionMethod: 0, filterMethod: 0, interlaceMethod: 0},
			},
			expectedResult: 1,
		},
		{
			testName: "color type 2 - Each pixel is an R,G,B triple",
			png: pngImage{
				meta: imageMetadata{width: 32, height: 32, bitDepth: 8, colorType: 2, compressionMethod: 0, filterMethod: 0, interlaceMethod: 0},
			},
			expectedResult: 3,
		},
		{
			testName: "color type 3 - Each pixel is a palette index",
			png: pngImage{
				meta: imageMetadata{width: 32, height: 32, bitDepth: 8, colorType: 3, compressionMethod: 0, filterMethod: 0, interlaceMethod: 0},
			},
			expectedResult: 1,
		},
		{
			testName: "color type 4 - Each pixel is a grayscale sample, followed by an alpha sample",
			png: pngImage{
				meta: imageMetadata{width: 32, height: 32, bitDepth: 8, colorType: 4, compressionMethod: 0, filterMethod: 0, interlaceMethod: 0},
			},
			expectedResult: 2,
		},
		{
			testName: "color type 6 -Each pixel is an R,G,B triple, followed by an alpha sample",
			png: pngImage{
				meta: imageMetadata{width: 32, height: 32, bitDepth: 16, colorType: 6, compressionMethod: 0, filterMethod: 0, interlaceMethod: 0},
			},
			expectedResult: 4,
		},
	}

	for _, testCase := range testTable {
		samplesPerPixel := testCase.png.samplesPerPixel()
		c.Assert(samplesPerPixel, Equals, testCase.expectedResult)
	}
}

func (s *MySuitePNG) TestStride(c *C) {

	testTable := []struct {
		testName       string
		png            pngImage
		expectedResult uint32
	}{
		{
			testName: "should pass - one byte per pixel",
			png: pngImage{
				meta: imageMetadata{width: 32, height: 32, bitDepth: 1, colorType: 0, compressionMethod: 0, filterMethod: 0, interlaceMethod: 0},
			},
			expectedResult: 32,
		},
		{
			testName: "should pass - 3 bytes per pixel",
			png: pngImage{
				meta: imageMetadata{width: 32, height: 32, bitDepth: 8, colorType: 2, compressionMethod: 0, filterMethod: 0, interlaceMethod: 0},
			},
			expectedResult: 96,
		},
		{
			testName: "should pass - 6 bytes per pixel",
			png: pngImage{
				meta: imageMetadata{width: 32, height: 32, bitDepth: 16, colorType: 2, compressionMethod: 0, filterMethod: 0, interlaceMethod: 0},
			},
			expectedResult: 192,
		},
	}
	for _, testCase := range testTable {
		stride := testCase.png.stride()
		c.Assert(stride, DeepEquals, testCase.expectedResult)
	}
}

func (s *MySuitePNG) TestUnfilter(c *C) {
	myPNG := pngImage{
		meta: imageMetadata{width: 8, bitDepth: 8, colorType: Grayscale, height: 4},
	}

	data := []byte{0x0, 0xff, 0x0, 0x8, 0xff, 0x8, 0xf, 0xff, 0x10, 0x0, 0x17, 0xff, 0x18, 0x1f, 0xff, 0x20, 0x27, 0xff, 0x0, 0x29, 0x2f, 0xff, 0x31, 0x37, 0xff, 0x39, 0x3f, 0x0, 0xff, 0x41, 0x47, 0xff, 0x4a, 0x4f, 0xff, 0x52}
	expected := []byte{0xff, 0x0, 0x8, 0xff, 0x8, 0xf, 0xff, 0x10, 0x17, 0xff, 0x18, 0x1f, 0xff, 0x20, 0x27, 0xff, 0x29, 0x2f, 0xff, 0x31, 0x37, 0xff, 0x39, 0x3f, 0xff, 0x41, 0x47, 0xff, 0x4a, 0x4f, 0xff, 0x52}
	err := myPNG.Unfilter(data)
	c.Assert(err, IsNil)
	c.Assert(myPNG.data, DeepEquals, expected)
}

func (s *MySuitePNG) TestProcessData(c *C) {

	pngJSON, err := ioutil.ReadFile("../../../test_files/png_test_files/expected/expected_parsed_json/01.basic_formats/basn0g01.txt")
	c.Assert(err, IsNil)
	myPNGstruct := new(png.StructPNG)
	err = json.Unmarshal(pngJSON, myPNGstruct)
	c.Assert(err, IsNil)

	expected := []byte{0xff, 0xff, 0xff, 0xfe, 0xff, 0xff, 0xff, 0xfc, 0xff, 0xff, 0xff, 0xf8, 0xff, 0xff, 0xff, 0xf0, 0xf3, 0xf3, 0xff, 0xe0, 0xf3, 0xf3, 0xff, 0xc0, 0xf3, 0xf3, 0xff, 0x80, 0xf3, 0x33, 0xff, 0x0, 0xf3, 0x33, 0xfe, 0x0, 0xf3,
		0x33, 0xfc, 0x0, 0xf8, 0x7, 0xf8, 0x0, 0xf8, 0x7, 0xf0, 0x0, 0xfc, 0xcf, 0xe0, 0x0, 0xfc, 0xcf, 0xc0, 0x0, 0xff, 0xff, 0x80,
		0x0, 0xff, 0xff, 0x0, 0x0, 0xff, 0xfe, 0x0, 0x0, 0xff, 0xfc, 0x0, 0x0, 0xff, 0xf8, 0xf, 0xe0, 0xff, 0xf0, 0xf, 0xe0, 0xff, 0xe0, 0xc, 0x30, 0xff, 0xc0, 0xc, 0x30, 0xff, 0x80, 0xf, 0xe0, 0xff, 0x0, 0xf, 0xe0, 0xfe, 0x0, 0xc, 0x30, 0xfc, 0x0, 0xc, 0x30, 0xf8, 0x0, 0xf, 0xe0, 0xf0, 0x0, 0xf, 0xe0, 0xe0, 0x0, 0x0, 0x0, 0xc0, 0x0, 0x0, 0x0, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	err = s.myPNG.ProcessData(myPNGstruct)
	c.Assert(err, IsNil)
	c.Assert(s.myPNG.data, DeepEquals, expected)

}

func (s *MySuitePNG) TestProcessImage(c *C) {

	// should Pass
	_, err := ProcessImage("../../../test_files/png_test_files/test_files/01.basic_formats/basn0g01.png")
	c.Assert(err, IsNil)

	// should Fail 1
	_, err = ProcessImage("../../../test_files/png_test_files/expected/expected_parsed_json/01.basic_formats/basn0g01")
	c.Assert(err, NotNil)

	// should Fail 2
	_, err = ProcessImage("../../../test_files/png_test_files/expected/expected_parsed_json/01.basic_formats/basn0g01.txt")
	c.Assert(err, Equals, png.ErrNotPNG)

}
