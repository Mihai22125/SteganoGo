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
