package pngint

import (
	"math"
)

// recon struct that holds reconstructed data
type Filterer struct {
	recon         []uint8 // a slice of reconstructed data
	stride        uint8
	height        uint32
	bytesPerPixel uint8
	bitDepth      uint8
}

// newRecon initialise new recon object
func newFilterer(bytesPerIndex uint8, stride uint8, height uint32, bitDepth uint8) *Filterer {
	f := new(Filterer)
	f.bytesPerPixel = bytesPerIndex
	f.stride = stride
	f.height = height
	f.bitDepth = bitDepth
	return f
}

// PaethPredicator algorithm used for Paeth filtering type
func (f *Filterer) PaethPredicator(a, b, c uint8) uint8 {
	var p float64
	var pr uint8

	p = float64(a) + float64(b) - float64(c)
	pa := math.Abs(p - float64(a))
	pb := math.Abs(p - float64(b))
	pc := math.Abs(p - float64(c))

	if pa <= pb && pa <= pc {
		pr = a
	} else if pb <= pc {
		pr = b
	} else {
		pr = c
	}

	return pr
}

// recon A reconstructs A, the byte corresponding to x in the pixel immediately before the pixel containing x
func (f *Filterer) reconA(scanLine uint32, byteIndex uint8) uint8 {
	if byteIndex >= f.bytesPerPixel {
		return f.recon[scanLine*uint32(f.stride)+uint32(byteIndex-f.bytesPerPixel)]
	}
	return 0
}

//recon B reconstructs B, the byte corresponding to x in the previous scanline
func (f *Filterer) reconB(scanLine uint32, byteIndex uint8) uint8 {
	if scanLine > 0 {
		return f.recon[(scanLine-1)*uint32(f.stride)+uint32(byteIndex)]
	}
	return 0
}

// recon C reconstructs C, the byte corresponding to b in the pixel immediately before the pixel containing b
func (f *Filterer) reconC(scanLine uint32, byteIndex uint8) uint8 {
	if scanLine > 0 && byteIndex >= f.bytesPerPixel {
		return f.recon[(scanLine-1)*uint32(f.stride)+uint32(byteIndex-f.bytesPerPixel)]
	}
	return 0
}

// reconstruct defiltres decompressed png data
func (f *Filterer) reconstruct(IDATdata []byte) error {
	i := 0
	f.recon = []byte{}

	if len(IDATdata) == 0 {
		return ErrInvalidInput
	}

	for row := uint32(0); row < f.height; row++ { // for each scanline
		filterType := IDATdata[i] // first byte of scanline is filter type
		i++
		for c := uint8(0); c < uint8(float64(f.stride)*(float64(f.bitDepth)/8)); c++ { // for each byte in scanline
			filtX := IDATdata[i]
			reconX := uint8(0)
			i++
			switch filterType {
			case byte(FiltNone):
				reconX = filtX
			case byte(FiltSub):
				reconX = filtX + f.reconA(row, c)
			case byte(FiltUp):
				reconX = filtX + f.reconB(row, c)
			case byte(FiltAverage):
				reconX = filtX + uint8((int(f.reconA(row, c))+int(f.reconB(row, c)))/2)
			case byte(FiltPaeth):
				reconX = filtX + f.PaethPredicator(f.reconA(row, c), f.reconB(row, c), f.reconC(row, c))
			default:
				return ErrUnknownFilterType
			}
			f.recon = append(f.recon, reconX)

		}
	}
	return nil
}

// The absolute value of a byte interpreted as a signed int8.
func abs8(d uint8) int {

	if d < 128 {
		return int(d)
	}
	return 256 - int(d)
}

// TODO: implement all filtering methods
// FilterData filter given byte slice
func (f *Filterer) FilterData(data []byte) []byte {

	size := int(f.stride)
	pr := make([]uint8, size)
	filtered := []byte{}
	// cr[*] and pr are the bytes for the current and previous row.
	// cr[0] is unfiltered (or equivalently, filtered with the ftNone filter).
	// cr[ft], for non-zero filter types ft, are buffers for transforming cr[0] under the
	// other PNG filter types. These buffers are allocated once and re-used for each row.
	// The +1 is for the per-row filter type, which is at cr[*][0].

	_ = size
	_ = pr

	cr := make([][]uint8, 6)
	for i := 0; i < 6; i++ {
		cr[i] = make([]uint8, size) // initialize a slice of dx unit8 in each of dy slices
	}

	for i := 0; i < len(data)/size; i++ {
		filtered = append(filtered, 0)
		filtered = append(filtered, data[i*size:i*size+size]...)
	}

	return filtered

}
