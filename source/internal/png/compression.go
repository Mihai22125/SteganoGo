package pngint

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
)

type Compressor struct {
}

func NewCompressor() *Compressor {
	return new(Compressor)
}

// DecompressPNGData returns decompressed data based on given method
func (c *Compressor) DecompressPNGData(data []byte, method CompressionMethod) ([]byte, error) {

	decompressed := []byte{}
	var err error

	if method != ComprDeflate {
		return nil, ErrNotSupportedPNG
	}

	if method == ComprDeflate {
		decompressed, err = c.deflate(data)
		if err != nil {
			return nil, err
		}
	}

	return decompressed, nil
}

// deflate : method 0 - DEFLATE/INFLATE compression
func (c *Compressor) deflate(data []byte) ([]byte, error) {

	reader := bytes.NewReader(data)
	decompressedReader, err := zlib.NewReader(reader)
	if err != nil {
		return nil, err
	}

	decompressed, err := ioutil.ReadAll(decompressedReader)
	if err != nil {
		return nil, err
	}
	decompressedReader.Close()

	return decompressed, nil
}

// CompressPNGData returns compressed data based on given method
func (c *Compressor) CompressPNGData(data []byte, method CompressionMethod) ([]byte, error) {

	if method != ComprDeflate {
		return nil, ErrNotSupportedPNG
	}

	compressed, err := c.inflate(data)
	if err != nil {
		return nil, err
	}

	return compressed, nil
}

// inflate : method 0 - DEFLATE/INFLATE compression
func (c *Compressor) inflate(data []byte) ([]byte, error) {

	var reader bytes.Buffer

	compressedReader := zlib.NewWriter(&reader)
	compressedReader.Write(data)
	compressedReader.Close()
	compressed := reader.Bytes()
	compressedReader.Close()

	return compressed, nil
}
