package pngint

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
)

// DecompressPNGData returns decompressed data based on given method
func DecompressPNGData(data []byte, method uint8) ([]byte, error) {

	decompressed := []byte{}
	var err error

	if method != ComprDeflate {
		return nil, ErrNotSupportedPNG
	}

	if method == ComprDeflate {
		decompressed, err = deflate(data)
		if err != nil {
			return nil, err
		}
	}

	return decompressed, nil
}

// deflate : method 0 - DEFLATE/INFLATE compression
func deflate(data []byte) ([]byte, error) {
	decompressed := []byte{}

	reader := bytes.NewReader(data)
	decompressedReader, err := zlib.NewReader(reader)
	if err != nil {
		return nil, err
	}

	decompressed, err = ioutil.ReadAll(decompressedReader)
	if err != nil {
		return nil, err
	}
	decompressedReader.Close()

	return decompressed, nil
}

// CompressPNGData returns compressed data based on given method
func CompressPNGData(data []byte, method uint8) ([]byte, error) {

	compressed := []byte{}

	if method != ComprDeflate {
		return nil, ErrNotSupportedPNG
	}

	if method == ComprDeflate {

	}

	return compressed, nil
}

// deflate : method 0 - DEFLATE/INFLATE compression
func inflate(data []byte) ([]byte, error) {
	compressed := []byte{}
	var reader bytes.Buffer

	compressedReader := zlib.NewWriter(&reader)
	compressedReader.Write(data)
	compressed = reader.Bytes()
	compressedReader.Close()

	return compressed, nil
}
