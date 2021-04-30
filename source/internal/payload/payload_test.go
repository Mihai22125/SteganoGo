package payload_test

import (
	"bytes"
	"testing"

	"github.com/Mihai22125/SteganoGo/internal/payload"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.

func Test(t *testing.T) { TestingT(t) }

type MySuite struct {
}

var _ = Suite(&MySuite{})

func (s *MySuite) SetUpSuite(c *C) {
}

type fakeImg struct {
	data     []byte
	bitDepth int
}

func (f *fakeImg) GetData() ([]byte, error) {
	return f.data, nil
}

func (f *fakeImg) UpdateData(data []byte) error {
	f.data = data
	return nil
}

func (f *fakeImg) GetBitDepth() int {
	return f.bitDepth
}

func (f *fakeImg) ReconstructImage(filepath string) error {
	return nil
}

func (s *MySuite) TestNewPayload(c *C) {
	data := []byte{1, 2, 3, 4, 5}
	ext := "ext"

	buf := bytes.NewBuffer(data)

	myPayload, err := payload.NewPayload(buf, ext)
	c.Assert(err, IsNil)

	c.Assert(myPayload.Data(), DeepEquals, data)
	c.Assert(myPayload.FileExtension(), Equals, ext)
	c.Assert(myPayload.ExtSize(), Equals, uint16(len(ext)))
	c.Assert(myPayload.Size(), Equals, uint32(len(data)))
}

func (s *MySuite) TestExtractPayload(c *C) {
	tt := []struct {
		data                []byte
		expectedPayloadData []byte
		expectedExt         string
	}{
		{
			data:                []byte{2, 0, 0, 0, 3, 0, 101, 120, 116, 121, 78},
			expectedPayloadData: []byte{121, 78},
			expectedExt:         "ext",
		},
	}

	for _, tc := range tt {
		buf := bytes.NewBuffer(tc.data)
		myPayload := payload.ExtractPayload(buf)
		_ = myPayload

		c.Assert(myPayload.Data(), DeepEquals, tc.expectedPayloadData)
		c.Assert(myPayload.FileExtension(), Equals, tc.expectedExt)
		c.Assert(myPayload.ExtSize(), Equals, uint16(len(tc.expectedExt)))
		c.Assert(myPayload.Size(), Equals, uint32(len(tc.expectedPayloadData)))
	}
}

func (s *MySuite) TestInsertPayload(c *C) {
	tt := []struct {
		data        []byte
		ext         string
		myFakeImg   *fakeImg
		expected    []byte
		expectedErr error
	}{
		{
			data: []byte{121, 78},
			ext:  "ext",
			myFakeImg: &fakeImg{
				data: []byte{
					20, 20, 20, 20, 20, 20, 20, 20, 31, 31, 31, 31, 31, 31, 31, 31, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1,
					0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1,
					0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 1, 1, 1,
				},
				bitDepth: 8,
			},
			expected: []byte{0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x15, 0x14, 0x1e, 0x1e, 0x1e, 0x1e, 0x1e, 0x1e, 0x1e, 0x1e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x1, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x1,
				0x1, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x1, 0x1, 0x1, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x1, 0x1, 0x1, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x0, 0x1, 0x1, 0x1, 0x0, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1},
			expectedErr: nil,
		},
	}

	for _, tc := range tt {

		buf := bytes.NewBuffer(tc.data)
		myPayload, err := payload.NewPayload(buf, tc.ext)
		c.Assert(err, IsNil)

		err = myPayload.InsertPayload(tc.myFakeImg)
		c.Assert(err, Equals, tc.expectedErr)

		modifiedData, err := tc.myFakeImg.GetData()
		c.Assert(err, IsNil)
		c.Assert(modifiedData, DeepEquals, tc.expected)
	}
}

func (s *MySuite) TestParsePayload(c *C) {
	tt := []struct {
		myImg          *fakeImg
		expectedResult []byte
		expectedErr    error
	}{
		{
			myImg: &fakeImg{
				data: []byte{0x14, 0x14, 0x14, 0x14, 0x14, 0x14, 0x15, 0x14, 0x1e, 0x1e, 0x1e, 0x1e, 0x1e, 0x1e, 0x1e, 0x1e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x1, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x1,
					0x1, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x1, 0x1, 0x1, 0x0, 0x1, 0x0, 0x0, 0x0, 0x1, 0x1, 0x1, 0x1, 0x0, 0x0, 0x1, 0x0, 0x1, 0x0, 0x0, 0x1, 0x1, 0x1, 0x0, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1},
				bitDepth: 8,
			},
			expectedResult: []byte{0x2, 0x0, 0x0, 0x0, 0x3, 0x0, 0x65, 0x78, 0x74, 0x79, 0x4e},
			expectedErr:    nil,
		},
	}

	for _, tc := range tt {
		result, err := payload.ParsePayload(tc.myImg)
		c.Assert(err, Equals, tc.expectedErr)
		c.Assert(result, DeepEquals, tc.expectedResult)
	}
}
