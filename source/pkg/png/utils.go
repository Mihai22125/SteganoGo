package png

// converts uint32 value to byte slice
func i32ToB(val uint32) []byte {
	r := make([]byte, 4)
	for i := uint32(0); i < 4; i++ {
		r[3-i] = byte((val >> (8 * i)) & 0xff)
	}
	return r
}
