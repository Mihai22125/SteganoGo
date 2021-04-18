package png

import (
	"bytes"
	"encoding/binary"
)

// Encoder configures encoding PNG images
type encoder struct {
}

// writeChunk writes a single Chunk to a bytes.Buffer
func (e *encoder) writeChunk(buf *bytes.Buffer, chunk Chunk) error {

	err := binary.Write(buf, binary.BigEndian, chunk.size)
	if err != nil {
		return err
	}

	err = binary.Write(buf, binary.BigEndian, []byte(chunk.chunkType))
	if err != nil {
		return err
	}

	err = binary.Write(buf, binary.BigEndian, chunk.data)
	if err != nil {
		return err
	}

	err = binary.Write(buf, binary.BigEndian, chunk.crc)
	if err != nil {
		return err
	}

	return nil
}

// writeChunks writes all chunks from a pngStruct to a bytes.Buffer
func (e *encoder) writeChunks(buf *bytes.Buffer, image *StructPNG) error {

	for _, chunk := range image.Chunks() {
		err := e.writeChunk(buf, chunk)
		if err != nil {
			return err
		}
	}

	return nil
}

// Encode encodes a pngStruct to a bytes Buffer
func (png *StructPNG) Encode() (*bytes.Buffer, error) {
	e := new(encoder)
	buf := new(bytes.Buffer)

	// write png header
	err := binary.Write(buf, binary.BigEndian, pngHeader)
	if err != nil {
		return buf, err
	}

	err = e.writeChunks(buf, png)

	return buf, err

}
