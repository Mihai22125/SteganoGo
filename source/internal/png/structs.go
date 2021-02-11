package pngint

type imageMetadata struct {
	width             uint32
	height            uint32    // gives the image dimensions in pixels. Zero is an invalid value.
	bitDepth          uint8     // gives the number of bits per sample or per palette index (not per pixel)
	colorType         ColorType // defines the PNG image type
	compressionMethod uint8     // indicates the method used to compress the image data
	filterMethod      uint8     // indicates the preprocessing method applied to the image data before compression
	interlaceMethod   uint8     // indicates whether there is interlacing
}

// recon slice that holds reconstructed data
type recon []uint8

// pngImage structure with parsed png data
type pngImage struct {
	meta imageMetadata
	data recon
}
