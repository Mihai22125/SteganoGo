package png

// Header : a PNG file starts with an 8-byte signature
type Header struct {
	header uint64
}

// Chunk conveys certain information about the image
type Chunk struct {
	size      uint32
	chunkType uint32
	data      []byte
	crc       uint32
}

// StructPNG stores header and chunks of PNG file
type StructPNG struct {
	header Header
	chunks []Chunk
}

// GetPNGHeader return png header
func (p StructPNG) GetPNGHeader() Header {
	return p.header
}

// GetPNGChunks return a slice of chunks
func (p StructPNG) GetPNGChunks() []Chunk {
	return p.chunks
}

// GetChunkSize return chunk size
func (c Chunk) GetChunkSize() uint32 {
	return c.size
}

// GetChunkType return chunk type
func (c Chunk) GetChunkType() uint32 {
	return c.chunkType
}

// GetChunkData return chunk data
func (c Chunk) GetChunkData() []byte {
	return c.data
}

// GetChunkCRC return chunk CRC
