package img

type Img interface {

	// Methods

	// GetData returns image's pixels data in a slice of bytes
	GetData() ([]byte, error)

	// UpdateData updates image's pixels data
	UpdateData(newImageData []byte) error

	// GetBitDepth returns the number of bits per sample or per palette index (not per pixel). Valid values are 1, 2, 4, 8, and 16.
	GetBitDepth() int

	// ReconstructImage reconstructs the image and saves it in the given path
	ReconstructImage(filepath string) error
}
