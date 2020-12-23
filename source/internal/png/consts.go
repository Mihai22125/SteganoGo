package pngint

// png compression method
// ComprDeflate compression method 0
// deflate/inflate compression with a sliding window of at most 32768 bytes
const ComprDeflate uint8 = 0

// png filter method
const FilterAdaptive uint8 = 0 // adaptive filtering with five basic filter types

// png interlace method
const NoInterlace uint8 = 0    // no interlace
const InterlaceAdam7 uint8 = 1 // Adam7 interlace

// png image types
const GrayscalePNG uint8 = 0          // each pixel is a grayscale sample
const TruecolorPNG uint8 = 2          // each pixel is a R,G,B triple
const IndexedColorPNG uint8 = 3       // each pixel is a palette index; a PLTE chunk shall appear
const GrayscaleWithAlphaPNG uint8 = 4 // each pixel is a grayscale sample followed by an alpha sample
const TruecolorWithAlphaPNG uint8 = 6 // each pixel is a R,G,B triple followed by an alpha sample
