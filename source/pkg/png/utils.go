package png

import "bytes"

// converts uint32 value to byte slice
func i32ToB(val uint32) []byte {
	r := make([]byte, 4)
	for i := uint32(0); i < 4; i++ {
		r[3-i] = byte((val >> (8 * i)) & 0xff)
	}
	return r
}

// CompareType compares uint32 type value with another defined png type
func CompareType(val uint32, pngType []byte) bool {
	if bytes.Equal(i32ToB(val), pngType) {
		return true
	}
	return false
}
