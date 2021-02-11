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
