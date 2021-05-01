package payload

type Payload struct {
	header header // payload header
	data   []byte // payload data
}

type header struct {
	size          uint32 // size of payload, without header
	extSize       uint16 // file extension size
	fileExtension string // output file extension
}

func (p *Payload) Size() uint32 {
	return p.header.size
}

func (p *Payload) ExtSize() uint16 {
	return p.header.extSize
}

func (p *Payload) FileExtension() string {
	return p.header.fileExtension
}

func (p *Payload) Data() []byte {
	return p.data
}
