package pngint

// CompressionMethod png compression method
type CompressionMethod uint8

// InterlaceMethod png interlace method
type InterlaceMethod uint8

// ColorType png color type
type ColorType uint8

// FiltMethod IDAT chunk filtering method
type FiltMethod uint8

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
const Grayscale ColorType = 0

// Truecolor each pixel is a R,G,B triple
const Truecolor ColorType = 2

// IndexedColor each pixel is a palette index; a PLTE chunk shall appear
const IndexedColor ColorType = 3

// GrayscaleWithAlpha each pixel is a grayscale sample followed by an alpha sample
const GrayscaleWithAlpha ColorType = 4

// TruecolorWithAlpha each pixel is a R,G,B triple followed by an alpha sample
const TruecolorWithAlpha ColorType = 6

// FiltNone Filter function: Filt(x) = Orig(x)
// Reconstruction function: Recon(x) = Filt(x)
const FiltNone FiltMethod = 0

// FiltSub Filter function: Filt(x) = Orig(x) - Orig(a)
// Reconstruction function: Recon(x) = Filt(x) + Recon(a)
const FiltSub FiltMethod = 1

// FiltUp Filter function: Filt(x) = Orig(x) - Orig(b)
// Reconstruction function: Recon(x) = Filt(x) + Recon(b)
const FiltUp FiltMethod = 2

// FiltAverage Filter function: Filt(x) = Orig(x) - floor((Orig(a) + Orig(b)) / 2)
// Reconstruction function: Recon(x) = Filt(x) + floor((Recon(a) + Recon(b)) / 2)
const FiltAverage FiltMethod = 3

// FiltPaeth Filter function: Filt(x) = Orig(x) - PaethPredictor(Orig(a), Orig(b), Orig(c))
// Reconstruction function: Recon(x) = Filt(x) + floor((Recon(a) + Recon(b)) / 2)
const FiltPaeth FiltMethod = 4

/*
where:

    x is the byte being filtered
    a is the byte corresponding to x in the pixel immediately before the pixel containing x
    b is the byte corresponding to x in the previous scanline
	c is the byte corresponding to b in the pixel immediately before the pixel containing b
*/
