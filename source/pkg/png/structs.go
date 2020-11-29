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
