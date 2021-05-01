package pngint

import (
	"hash/crc32"
	"os"

	"github.com/Mihai22125/SteganoGo/pkg/png"
)

func (pngImg *PngImage) reconstructIDAT() ([]byte, error) {

	filterer := newFilterer(pngImg.bytesPerPixel(), uint16(pngImg.stride()), pngImg.meta.height, pngImg.meta.bitDepth)
	compressor := new(Compressor)

	filteredData := filterer.FilterData(pngImg.data)
	compressed, err := compressor.CompressPNGData(filteredData, ComprDeflate)

	if err != nil {
		return []byte{}, err
	}

	return compressed, nil
}

func (pngImg *PngImage) divideIDATChunks(data []byte) []png.Chunk {
	var divided [][]byte
	var idatChunks []png.Chunk

	chunkSize := 8192

	for i := 0; i < len(data); i += chunkSize {
		end := i + chunkSize

		if end > len(data) {
			end = len(data)
		}

		divided = append(divided, data[i:end])
	}

	for _, ch := range divided {
		crc := crc32.ChecksumIEEE(append([]byte(png.TypeIDAT), ch...))
		chunk := png.NewChunk(uint32(len(ch)), png.TypeIDAT, ch, crc)
		idatChunks = append(idatChunks, chunk)
	}

	return idatChunks
}

func (pngImg *PngImage) UpdatePNG() error {

	pngData, err := pngImg.reconstructIDAT()
	if err != nil {
		return err
	}

	idatChunks := pngImg.divideIDATChunks(pngData)
	pngImg.png.UpdateIdatChunks(idatChunks)

	return nil
}

func (pngImg *PngImage) ReconstructImage(filepath string) error {

	buf, err := pngImg.png.Encode()
	if err != nil {
		return err
	}

	// open output file
	fo, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}

	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	fo.Write(buf.Bytes())

	return nil
}
