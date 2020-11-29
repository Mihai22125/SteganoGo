package png

import (
	"SteganoGo/pkg/fileutils"
	"encoding/binary"
	"fmt"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	path := "C:\\Users\\Mihai\\Desktop\\6_1_ss.png"

	fHandle, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		//return
	}

	processedFile, err := fileutils.PreProcessFile(fHandle)
	if err != nil {
		fmt.Println(err)
		//return
	}

	parsed, err := ParsePNG(processedFile)
	if err != nil {
		fmt.Println(err)
		//return
	}

	fmt.Println(parsed.header)
	for i, chunk := range parsed.chunks {
		fmt.Println(i)
		fmt.Println("size: ", chunk.size)
		b := make([]byte, 8)
		binary.BigEndian.PutUint64(b, uint64(chunk.chunkType))

		fmt.Println("type: ", string(b))
		//fmt.Println("data: ", chunk.data)
		fmt.Println("CRC: ", chunk.crc)
		fmt.Println()

	}

	//fmt.Println(parsed)

}
