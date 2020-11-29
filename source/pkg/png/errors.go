package png

import "errors"

// ErrBadPNG ..
var ErrBadPNG error = errors.New("PNG file format is not correct")

// ErrNotPNG ..
var ErrNotPNG error = errors.New("Given file is not a PNG file")
