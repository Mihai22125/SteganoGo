package pngint

// CompressionMethod png compression method
type CompressionMethod uint8

// InterlaceMethod png interlace method
type InterlaceMethod uint8

// ImgType png image types
type ImgType uint8

// ComprDeflate compression method 0
// deflate/inflate compression with a sliding window of at most 32768 bytes
const ComprDeflate CompressionMethod = 0

// FilterAdaptive png filter method
const FilterAdaptive CompressionMethod = 0 // adaptive filtering with five basic filter types

// NoInterlace no interlace
const NoInterlace InterlaceMethod = 0

// InterlaceAdam7 Adam7 interlace
const InterlaceAdam7 InterlaceMethod = 1

// Grayscale each pixel is a grayscale sample
const Grayscale ImgType = 0

// Truecolor each pixel is a R,G,B triple
const Truecolor ImgType = 2

// IndexedColor each pixel is a palette index; a PLTE chunk shall appear
const IndexedColor ImgType = 3

// GrayscaleWithAlpha each pixel is a grayscale sample followed by an alpha sample
const GrayscaleWithAlpha ImgType = 4

// TruecolorWithAlpha each pixel is a R,G,B triple followed by an alpha sample
const TruecolorWithAlpha ImgType = 6
