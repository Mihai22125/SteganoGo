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

// Header return png header
func (p StructPNG) Header() Header {
	return p.header
}

// Chunks return a slice of chunks
func (p StructPNG) Chunks() []Chunk {
	return p.chunks
}

// Size return chunk size
func (c Chunk) Size() uint32 {
	return c.size
}

// ChunkType return chunk type
func (c Chunk) ChunkType() uint32 {
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
