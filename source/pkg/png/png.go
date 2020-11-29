package png

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// package with functions used for parsing PNG files

// ParsePNG extracts data from PNG file into a struct
func ParsePNG(r *bytes.Reader) (StructPNG, error) {
	newPNG := StructPNG{}

	buf := make([]byte, 8)

	// read file header
	len, err := r.Read(buf)
	if err != nil {
		fmt.Println("[ParsePNG]: failed to read PNG file: ", err)
		return newPNG, err
	}
	if len != 8 {
		return newPNG, ErrBadPNG
	}
	// check if file is indeed an PNG file
	if checkPNG(buf) == false {
		return newPNG, ErrNotPNG
	}

	newPNG.chunks, err = readChunks(r)
	if err != nil {
		return newPNG, err
	}

	return newPNG, nil
}

// checkPNG: returns true if the header is a known PNG header and false otherwise
func checkPNG(header []byte) bool {
	if bytes.Equal(header, pngHeader) {
		return true
	}
	return false
}

func readSingleChunck(r *bytes.Reader) (Chunk, error) {
	newChunk := Chunk{}

	// read chunk data length
	err := binary.Read(r, binary.BigEndian, &newChunk.size)
	if err != nil {
		fmt.Println("[ParsePNG]: failed to read PNG file: ", err)
		return newChunk, err
	}

	// read chunk type
	err = binary.Read(r, binary.BigEndian, &newChunk.chunkType)
	if err != nil {
		fmt.Println("[ParsePNG]: failed to read PNG file: ", err)
		return newChunk, err
	}

	newChunk.data = make([]byte, newChunk.size)
	// read chunk data
	err = binary.Read(r, binary.BigEndian, &newChunk.data)
	if err != nil {
		fmt.Println("[ParsePNG]: failed to read PNG file: ", err)
		return newChunk, err
	}
	// read chunk CRC
	err = binary.Read(r, binary.BigEndian, &newChunk.crc)
	if err != nil {
		fmt.Println("[ParsePNG]: failed to read PNG file: ", err)
		return newChunk, err
	}

	return newChunk, nil
}

func readChunks(r *bytes.Reader) ([]Chunk, error) {
	chunks := []Chunk{}

	for {
		newChunk, err := readSingleChunck(r)
		if err != nil {
			return chunks, err
		}
		chunks = append(chunks, newChunk)

		if bytes.Equal(i32ToB(newChunk.chunkType), typeIEND) {
			break
		}
	}

	return chunks, nil
}
