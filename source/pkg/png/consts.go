package png

// pngHeader : PNG file known header
// The first eight bytes of a PNG file always contain the following (decimal) values: 137 80 78 71 13 10 26 10
var pngHeader = []byte{137, 80, 78, 71, 13, 10, 26, 10}

// Critical chunk types

// TypeIHDR must be the first chunk
var TypeIHDR = "IHDR"

// TypePLTE contains the palette: a list of colors
var TypePLTE = "PLTE"

// TypeIDAT contains the image, which may be split among multiple IDAT chunks
// The IDAT chunk contains the actual image data, which is the output stream of the compression algorithm
var TypeIDAT = "IDAT"

// TypeIEND marks the image end; the data field of the IEND chunk has 0 bytes/is empty
var TypeIEND = "IEND"
