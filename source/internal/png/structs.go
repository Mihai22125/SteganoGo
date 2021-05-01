package pngint

import "github.com/Mihai22125/SteganoGo/pkg/png"

type imageMetadata struct {
	width             uint32
	height            uint32            // gives the image dimensions in pixels. Zero is an invalid value.
	bitDepth          uint8             // gives the number of bits per sample or per palette index (not per pixel)
	colorType         ColorType         // defines the PNG image type
	compressionMethod CompressionMethod // indicates the method used to compress the image data
	filterMethod      FiltMethod        // indicates the preprocessing method applied to the image data before compression
	interlaceMethod   InterlaceMethod   // indicates whether there is interlacing
}

// PngImage structure with parsed png data
type PngImage struct {
	meta imageMetadata
	data []uint8
	png  png.StructPNG
}

func (p *PngImage) GetData() ([]byte, error) {
	return p.data, nil
}

func (p *PngImage) UpdateData(newImageData []byte) error {
	p.data = newImageData
	return nil
}

func (p *PngImage) GetBitDepth() int {
	return int(p.meta.bitDepth)
}
