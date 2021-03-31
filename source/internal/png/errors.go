package pngint

import "errors"

// ErrBadPNG ..
var ErrBadPNG error = errors.New("PNG file format is not correct")

// ErrNotSupportedPNG ..
var ErrNotSupportedPNG error = errors.New("PNG file format is not supported")

// ErrNoIHDR ..
var ErrNoIHDR error = errors.New("IHDR chunk not found in PNG file")

// ErrBadSizeIHDR ..
var ErrBadSizeIHDR error = errors.New("IHDR chunk data size is not accepted")

// ErrBadIHDR ..
var ErrBadIHDR error = errors.New("IHDR chunk data is not accepted")

// ErrOutOfBoundsPixel given pixel is out of bounds
var ErrOutOfBoundsPixel error = errors.New("pixel is out of bounds")

// ErrUnknownFilterType unknown filter type
var ErrUnknownFilterType error = errors.New("unknown filter type")

// ErrInvalidInput invalid input passed to a specific function
var ErrInvalidInput error = errors.New("invalid input")
