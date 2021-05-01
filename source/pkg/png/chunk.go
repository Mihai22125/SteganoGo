package png

func NewChunk(size uint32, chunkType string, data []byte, crc uint32) Chunk {
	chunk := Chunk{}
	chunk.size = size
	chunk.chunkType = chunkType
	chunk.data = data
	chunk.crc = crc

	return chunk
}

// Size return chunk size
func (c Chunk) Size() uint32 {
	return c.size
}

// ChunkType return chunk type
func (c Chunk) ChunkType() string {
	return c.chunkType
}

// Data return chunk data
func (c Chunk) Data() []byte {
	return c.data
}

// CRC return chunk CRC
func (c Chunk) CRC() uint32 {
	return c.crc
}

// CompareType returns true if chunk type equals given type
func (ch Chunk) CompareType(chType string) bool {
	if ch.chunkType == chType {
		return true
	}
	return false
}
