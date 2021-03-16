package png

import "errors"

// ErrBadPNG ..
var ErrBadPNG error = errors.New("PNG file format is not correct")

// ErrNotPNG ..
var ErrNotPNG error = errors.New("Given file is not a PNG file")

var ErrPNGChunks error = errors.New("Failed to read PNG chunks")

// ErrIHDRMissing IHDR chunk is missing
var ErrIHDRMissing error = errors.New("IHDR chunk is missing")

// ErrIDATMissing IHDR chunk is missing
var ErrIDATMissing error = errors.New("IDAT chunk is missing")

// ErrInvalidIHDR invalid IHDR chunk
var ErrInvalidIHDR error = errors.New("IHDR chunk is not valid")
