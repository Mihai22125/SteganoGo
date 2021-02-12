package pngint

import (
	"encoding/binary"
	"fmt"

	"github.com/Mihai22125/SteganoGo/pkg/png"
)

// extractMetadata extracts data from IHDR chunk. Returns error
func (pngImg *pngImage) extractMetadata(stpng png.StructPNG) error {

	ihdr, err := stpng.IHDRChunk()
	if err != nil {
		fmt.Println("[extractMetadata]: failed to get IHDR chunk")
		return err
	}

	ihdrData := ihdr.Data()
	err = pngImg.processIHDR(ihdrData)
	if err != nil {
		return err
	}

	return nil
}

// parseIHDR perse IHDR chunk
func (pngImg *pngImage) processIHDR(ihdrData []byte) error {

	if len(ihdrData) != 13 {
		fmt.Println("IHDR chunk has invalid size")
		return png.ErrInvalidIHDR
	}
	meta := imageMetadata{}

	buf := ihdrData[0:4]
	meta.width = binary.LittleEndian.Uint32(buf)
	buf = ihdrData[4:8]
	meta.height = binary.LittleEndian.Uint32(buf)
	meta.bitDepth = ihdrData[8]
	meta.colorType = ColorType(ihdrData[9])
	meta.compressionMethod = ihdrData[10]
	meta.filterMethod = ihdrData[11]
	meta.interlaceMethod = ihdrData[12]

	pngImg.meta = meta
	return nil
}

// newRecon return an empty recon
func newRecon() recon {
	return recon{}
}

// bytesPerPixel retun bytes per pixel based on color type
func (pngImg *pngImage) bytesPerPixel() uint8 {
	if pngImg.meta.colorType == Grayscale || pngImg.meta.colorType == IndexedColor {
		return 1
	}
	if pngImg.meta.colorType == GrayscaleWithAlpha {
		return 2
	}
	if pngImg.meta.colorType == Truecolor {
		return 3
	}
	return 4
}

// stride return bytes per row from png image
func (pngImg *pngImage) stride() uint32 {
	return pngImg.meta.width * uint32(pngImg.bytesPerPixel())
}

// A is the byte corresponding to x in the pixel immediately before the pixel containing x
func (pngImg *pngImage) reconA(scanLine uint32, byteIndex uint8, bytesPerPixel uint8, stride uint32) (uint8, error) {
	if byteIndex >= bytesPerPixel {
		return pngImg.data[scanLine*stride+uint32(byteIndex-bytesPerPixel)], nil
	}
	return 0, ErrOutOfBoundsPixel
}

// B is the byte corresponding to x in the previous scanline
func (pngImg *pngImage) reconB(scanLine uint32, byteIndex uint8, stride uint32) (uint8, error) {
	if scanLine > 0 {
		return pngImg.data[(scanLine-1)*stride+uint32(byteIndex)], nil
	}
	return 0, ErrOutOfBoundsPixel
}

// C is the byte corresponding to b in the pixel immediately before the pixel containing b
func (pngImg *pngImage) reconC(scanLine uint32, byteIndex uint8, bytesPerPixel uint8, stride uint32) (uint8, error) {
	if scanLine > 0 && byteIndex >= bytesPerPixel {
		return pngImg.data[(scanLine-1)*stride+uint32(byteIndex-bytesPerPixel)], nil
	}
	return 0, ErrOutOfBoundsPixel
}