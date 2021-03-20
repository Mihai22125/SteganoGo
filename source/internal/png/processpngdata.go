package pngint

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"

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
	meta.compressionMethod = CompressionMethod(ihdrData[10])
	meta.filterMethod = FiltMethod(ihdrData[11])
	meta.interlaceMethod = InterlaceMethod(ihdrData[12])

	pngImg.meta = meta
	return nil
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

// Unfilter
func (pngImg *pngImage) Unfilter(decompressed []byte) error {

	filterer := newFilterer(pngImg.bytesPerPixel(), uint8(pngImg.stride()), pngImg.meta.height)

	// defilter uncompressed data
	err := filterer.reconstruct(decompressed)
	if err != nil {
		return err
	}

	// assign processed data to png struct
	pngImg.data = filterer.recon
	return nil
}

// ProcessData consumes an png.StructPNG and it processes png data
func (pngImg *pngImage) ProcessData(stpng *png.StructPNG) error {
	compressor := NewCompressor()

	IDATdata, err := stpng.IDATdata()
	if err != nil {
		return err
	}

	// decompress png data
	decompressed, err := compressor.DecompressPNGData(IDATdata, pngImg.meta.compressionMethod)
	if err != nil {
		return err
	}

	err = pngImg.Unfilter(decompressed)
	if err != nil {
		return err
	}

	return nil
}

func (pngImg *pngImage) ProcessImage(path string) {
	file, err := os.Open("file.go") // For read access.
	if err != nil {
		log.Fatal(err)
	}

	pngData, err := png.ParsePNG(file)
	if err != nil {
		log.Fatal(err)
	}

	newPNGImage := pngImage{}

	newPNGImage.ProcessData(&pngData)

}
