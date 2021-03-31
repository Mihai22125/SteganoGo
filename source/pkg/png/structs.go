package png

import (
	"encoding/json"
)

// Header : a PNG file starts with an 8-byte signature
type Header struct {
	signature []byte
}

// Chunk conveys certain information about the image
type Chunk struct {
	size      uint32
	chunkType string
	data      []byte
	crc       uint32
}

// StructPNG stores header and chunks of PNG file
type StructPNG struct {
	header Header
	chunks []Chunk
}

// Header return png header
func (p StructPNG) Header() Header {
	return p.header
}

// Chunks return a slice of chunks
func (p StructPNG) Chunks() []Chunk {
	return p.chunks
}

// MarshalJSON custom MarshalJSON for Header struct
func (h Header) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		Signature []byte `json:"signature"`
	}{
		Signature: h.signature,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (h *Header) UnmarshalJSON(b []byte) error {
	temp := &struct {
		Signature []byte `json:"signature"`
	}{
		Signature: h.signature,
	}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	h.signature = temp.Signature

	return nil
}

// MarshalJSON custom MarshalJSON for Header struct
func (c Chunk) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		Size      uint32 `json:"size"`
		ChunkType string `json:"chunk_type"`
		Data      []byte `json:"data,omitempty"`
		Crc       uint32 `json:"crc"`
	}{
		Size:      c.size,
		ChunkType: c.chunkType,
		Data:      c.data,
		Crc:       c.crc,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (c *Chunk) UnmarshalJSON(b []byte) error {
	temp := &struct {
		Size      uint32 `json:"size"`
		ChunkType string `json:"chunk_type"`
		Data      []byte `json:"data,omitempty"`
		Crc       uint32 `json:"crc"`
	}{
		Size:      c.size,
		ChunkType: c.chunkType,
		Data:      c.data,
		Crc:       c.crc,
	}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	c.chunkType = temp.ChunkType
	c.size = temp.Size
	c.data = temp.Data
	c.crc = temp.Crc

	return nil
}

// MarshalJSON custom MarshalJSON for Header struct
func (p StructPNG) MarshalJSON() ([]byte, error) {
	j, err := json.Marshal(struct {
		Header Header  `json:"header"`
		Chunks []Chunk `json:"chunks"`
	}{
		Header: p.header,
		Chunks: p.chunks,
	})
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (p *StructPNG) UnmarshalJSON(b []byte) error {
	temp := &struct {
		Header Header  `json:"header"`
		Chunks []Chunk `json:"chunks"`
	}{
		Header: p.header,
		Chunks: p.chunks,
	}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	p.header = temp.Header
	p.chunks = temp.Chunks

	return nil
}
