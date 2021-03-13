package png

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/Mihai22125/SteganoGo/pkg/fileutils"
)

// package with functions used for parsing PNG files

// ParsePNG consumes an os.File and extracts PNG raw data from it
func ParsePNG(dat *os.File) (StructPNG, error) {
	png := StructPNG{}
	bReader, err := fileutils.PreProcessFile(dat)
	if err != nil {
		return png, err
	}

	png, err = parsePNG(bReader)

	return png, err
}

// parsePNG consumes a bytes.Reader and returns PNG structure
func parsePNG(r *bytes.Reader) (StructPNG, error) {
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

	newPNG.header = Header{buf}

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

// readSingleChunck read the next chunk from png data
func readSingleChunck(r *bytes.Reader) (Chunk, error) {
	newChunk := Chunk{}

	// read chunk data length
	err := binary.Read(r, binary.BigEndian, &newChunk.size)
	if err != nil {
		fmt.Println("[readSingleChunck]: failed to read PNG chunk size: ", err)
		return newChunk, err
	}

	// read chunk type
	buf := make([]byte, 4)
	_, err = r.Read(buf)
	if err != nil {
		fmt.Println("[readSingleChunck]: failed to read PNG chunk type: ", err)
		return newChunk, err
	}
	newChunk.chunkType = string(buf)

	newChunk.data = make([]byte, newChunk.size)
	// read chunk data
	err = binary.Read(r, binary.BigEndian, &newChunk.data)
	if err != nil {
		fmt.Println("[readSingleChunck]: failed to read PNG chunk data: ", err)
		return newChunk, err
	}
	// read chunk CRC
	err = binary.Read(r, binary.BigEndian, &newChunk.crc)
	if err != nil {
		fmt.Println("[readSingleChunck]: failed to read PNG crc: ", err)
		return newChunk, err
	}

	return newChunk, nil
}

// readChunks read png chunks and returns a slice of Chunk
func readChunks(r *bytes.Reader) ([]Chunk, error) {
	chunks := []Chunk{}

	for {
		newChunk, err := readSingleChunck(r)
		if err != nil {
			return chunks, err
		}
		chunks = append(chunks, newChunk)
		if newChunk.chunkType == TypeIEND {
			break
		}
	}

	return chunks, nil
}

// CompareType returns true if chunk type equals given type
func (ch Chunk) CompareType(chType string) bool {
	if ch.chunkType == chType {
		return true
	}
	return false
}

// IHDRChunk returns png IHDR chunk
func (p StructPNG) IHDRChunk() (Chunk, error) {
	for _, chunk := range p.Chunks() {
		if chunk.CompareType(TypeIHDR) == true {
			return chunk, nil
		}
	}
	return Chunk{}, ErrIHDRMissing
}

// IDATChunks returns a slice of found IDAT chunks
func (p StructPNG) IDATChunks() ([]Chunk, error) {
	idat := []Chunk{}

	for _, chunk := range p.Chunks() {
		if chunk.CompareType(TypeIDAT) == true {
			idat = append(idat, chunk)
		}
	}
	if len(idat) == 0 {
		return nil, ErrIDATMissing
	}

	return idat, nil
}

// IDATdata returns IDAT chunks concatenated in a single slice of bytes
func (p *StructPNG) IDATdata() ([]byte, error) {
	data := []byte{}

	IDATChunks, err := p.IDATChunks()
	if err != nil {
		return nil, err
	}

	for _, chunk := range IDATChunks {
		data = append(data, chunk.data...)
	}

	return data, nil
}
