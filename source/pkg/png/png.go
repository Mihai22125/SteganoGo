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
	if err != nil || len != 8 {
		fmt.Println("[ParsePNG]: failed to read PNG file: ", err)
		return EmptyPNG, ErrBadPNG
	}

	newPNG.header = Header{buf}

	// check if file is indeed an PNG file
	if checkPNG(buf) == false {
		return EmptyPNG, ErrNotPNG
	}

	newPNG.chunks, err = readChunks(r)
	if err != nil {
		return EmptyPNG, ErrPNGChunks
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
		return Chunk{}, err
	}

	// read chunk type
	buf := make([]byte, 4)
	err = binary.Read(r, binary.BigEndian, buf)
	if err != nil {
		return Chunk{}, err
	}
	newChunk.chunkType = string(buf)

	newChunk.data = make([]byte, newChunk.size)

	// read chunk data
	err = binary.Read(r, binary.BigEndian, &newChunk.data)
	if err != nil {
		return Chunk{}, err
	}
	// read chunk CRC
	err = binary.Read(r, binary.BigEndian, &newChunk.crc)
	if err != nil {
		return Chunk{}, err
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

func (p *StructPNG) UpdateIdatChunks(newIDATChunks []Chunk) {
	var beforeIDAT, afterIDAT []Chunk
	f, _ := os.Create("before.txt")
	fmt.Fprintf(f, "%v+\n", p)
	i := 0

	for i = 0; i < len(p.chunks); i++ {
		if p.chunks[i].chunkType == TypeIDAT {
			beforeIDAT = p.chunks[:i]
			break
		}
	}

	for j := i; j < len(p.chunks); j++ {
		if p.chunks[j].chunkType != TypeIDAT {
			afterIDAT = p.chunks[j:]
		}
	}

	updatedChunks := []Chunk{}
	updatedChunks = append(updatedChunks, beforeIDAT...)
	updatedChunks = append(updatedChunks, newIDATChunks...)
	updatedChunks = append(updatedChunks, afterIDAT...)

	p.chunks = updatedChunks
	f, _ = os.Create("after.txt")
	fmt.Fprintf(f, "%v+\n", p)
}
