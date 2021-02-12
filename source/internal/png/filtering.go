package pngint

import (
	"math"
)

// PaethPredicator algorithm used for Paeth filtering type
func PaethPredicator(a, b, c uint8) uint8 {
	var p float64
	var pr uint8

	p = float64(a + b + c)
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

// newRecon initialise new recon object
func newRecon(bytesPerIndex uint8, stride uint8, height uint32) *recon {
	r := new(recon)
	r.bytesPerPixel = bytesPerIndex
	r.stride = stride
	r.height = height
	return r
}

// recon A reconstructs A, the byte corresponding to x in the pixel immediately before the pixel containing x
func (r *recon) reconA(scanLine uint32, byteIndex uint8) uint8 {
	if byteIndex >= r.bytesPerPixel {
		return r.recon[scanLine*uint32(r.stride)+uint32(byteIndex-r.bytesPerPixel)]
	}
	return 0
}

//recon B reconstructs B, the byte corresponding to x in the previous scanline
func (r *recon) reconB(scanLine uint32, byteIndex uint8) uint8 {
	if scanLine > 0 {
		return r.recon[(scanLine-1)*uint32(r.stride)+uint32(byteIndex)]
	}
	return 0
}

// recon C reconstructs C, the byte corresponding to b in the pixel immediately before the pixel containing b
func (r *recon) reconC(scanLine uint32, byteIndex uint8) uint8 {
	if scanLine > 0 && byteIndex >= r.bytesPerPixel {
		return r.recon[(scanLine-1)*uint32(r.stride)+uint32(byteIndex-r.bytesPerPixel)]
	}
	return 0
}

// reconstruct defiltres decompressed png data
func (r *recon) reconstruct(IDATdata []byte) error {
	i := 0
	for row := uint32(0); row < r.height; row++ { // for each scanline
		filterType := IDATdata[i] // first byte of scanline is filter type
		i++
		for c := uint8(0); c < r.stride; c++ { // for each byte in scanline
			filtX := IDATdata[i]
			reconX := uint8(0)
			i++
			switch filterType {
			case byte(FiltNone):
				reconX = filtX
			case byte(FiltSub):
				reconX = filtX + r.reconA(row, c)
			case byte(FiltUp):
				reconX = filtX + r.reconB(row, c)
			case byte(FiltAverage):
				reconX = filtX + (r.reconA(row, c)+r.reconB(row, c))/2
			case byte(FiltPaeth):
				reconX = filtX + PaethPredicator(r.reconA(row, c), r.reconB(row, c), r.reconC(row, c))
			default:
				return ErrUnknownFilterType
			}

			r.recon = append(r.recon, reconX)

		}
	}
	return nil
}
