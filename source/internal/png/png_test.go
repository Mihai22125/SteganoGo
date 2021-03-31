package pngint

import (
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
	/*
		testTable := []struct {
			name         string
			png          png.StructPNG
			expectedMeta imageMetadata
			expectedErr  error
		}{
			{
				name: "shouldPass1",

					png: png.StructPNG{

							chunks: []png.Chunk{
								{
									size:       13,
									chunk_type: png.TypeIHDR,
									data:       []byte{0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x20, 0x08, 0x00, 0x00, 0x00, 0x00},
									Crc:        1443964200,
								},
							},

				expectedMeta: imageMetadata{width: 32, height: 32, bitDepth: 8, colorType: 0, compressionMethod: 0, filterMethod: 0, interlaceMethod: 0},
				expectedErr:  nil,
			},


				{
					name: "shouldFail - invalid IHDR length",
					//IHDRchunk:    []byte{0x00, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00, 0x20, 0x08, 0x00, 0x00, 0x00},
					expectedMeta: imageMetadata{},
					expectedErr:  png.ErrInvalidIHDR,
				},

		}

		for _, testCase := range testTable {
			err := s.myPNG.extractMetadata(testCase.png)
			c.Assert(err, Equals, testCase.expectedErr)
			c.Assert(s.myPNG.meta, DeepEquals, testCase.expectedMeta)
		}
	*/
}
