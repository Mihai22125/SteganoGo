package payload

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"

	"github.com/Mihai22125/SteganoGo/internal/img"
)

func NewPayload(buf *bytes.Buffer, ext string) (Payload, error) {
	myPayload := Payload{}
	myPayload.header.fileExtension = ext

	data, err := ioutil.ReadAll(buf)
	if err != nil {
		return myPayload, err
	}

	myPayload.data = data
	myPayload.header.size = uint32(len(data))
	myPayload.header.extSize = uint16(len(ext))

	return myPayload, nil
}

func (p *Payload) GeneratePayload() *bytes.Buffer {
	buf := []byte{}

	buf = append(buf, []byte(PayloadSignature)...)
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, p.header.size)
	buf = append(buf, bs...)

	bs = make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, p.header.extSize)
	buf = append(buf, bs...)
	buf = append(buf, []byte(p.header.fileExtension)...)
	buf = append(buf, p.data...)
	return bytes.NewBuffer(buf)
}

func ExtractPayload(buf *bytes.Buffer) (Payload, error) {

	extractedPayload := Payload{}

	bs := make([]byte, len(PayloadSignature))
	buf.Read(bs)
	if string(bs) != PayloadSignature {
		return extractedPayload, ErrNotPayload
	}

	bs = make([]byte, 4)
	buf.Read(bs)
	extractedPayload.header.size = binary.LittleEndian.Uint32(bs)

	bs = make([]byte, 2)
	buf.Read(bs)
	extractedPayload.header.extSize = binary.LittleEndian.Uint16(bs)

	bs = make([]byte, extractedPayload.header.extSize)
	buf.Read(bs)
	extractedPayload.header.fileExtension = string(bs)

	bs = make([]byte, extractedPayload.header.size)
	buf.Read(bs)
	extractedPayload.data = bs

	return extractedPayload, nil
}

func (p *Payload) WriteFile() error {
	path := "extracted_payload." + p.header.fileExtension

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(p.data)
	if err != nil {
		return err
	}

	return nil
}

func (p *Payload) InsertPayload(myImage img.Img) error {

	imageData, err := myImage.GetData()
	if err != nil {
		return err
	}
	bitDepth := myImage.GetBitDepth()
	payloadDataBuf := p.GeneratePayload()
	payloadData, err := ioutil.ReadAll(payloadDataBuf)
	if err != nil {
		return err
	}

	availableBits := int(float64(len(imageData)) * (8. / float64(bitDepth)))
	payloadBitsSize := len(payloadData) * 8

	if bitDepth != 8 && bitDepth != 16 {
		return ErrUnsupportedBitDepth
	}
	if availableBits < payloadBitsSize {
		return ErrPayloadSize
	}

	k := 0
	for i := 0; i < len(payloadData); i++ {
		for j := 7; j >= 0; j-- {
			x := payloadData[i] >> j & 1
			if x == 1 {
				imageData[k] |= 1
			} else {
				imageData[k] &= ^byte(1)
			}
			k += bitDepth / 8
		}
	}

	myImage.UpdateData(imageData)
	return nil
}

func ParsePayload(myImage img.Img) ([]byte, error) {
	payload := []byte{}
	currentByte := byte(0)

	data, err := myImage.GetData()
	if err != nil {
		return payload, err
	}

	j := 7
	for i := 0; i < len(data); i++ {
		if j < 0 {
			payload = append(payload, currentByte)
			currentByte = 0
			j = 7
		}
		currentByte += (data[i] & 1) << j
		j--
	}
	return payload, nil
}
