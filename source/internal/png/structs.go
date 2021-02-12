package pngint

type imageMetadata struct {
	width             uint32
	height            uint32            // gives the image dimensions in pixels. Zero is an invalid value.
	bitDepth          uint8             // gives the number of bits per sample or per palette index (not per pixel)
	colorType         ColorType         // defines the PNG image type
	compressionMethod CompressionMethod // indicates the method used to compress the image data
	filterMethod      FiltMethod        // indicates the preprocessing method applied to the image data before compression
	interlaceMethod   InterlaceMethod   // indicates whether there is interlacing
}

//type recon []uint8

// recon struct that holds reconstructed data
type recon struct {
	recon         []uint8
	stride        uint8
	height        uint32
	bytesPerPixel uint8
}

// pngImage structure with parsed png data
type pngImage struct {
	meta imageMetadata
	data []uint8
}
