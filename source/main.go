package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Mihai22125/SteganoGo/internal/payload"
	pngint "github.com/Mihai22125/SteganoGo/internal/png"
)

func main() {
	decodePtr := flag.Bool("decode", false, "decode payload from an image.")
	insertPtr := flag.Bool("insert", false, "insert payload into an image.")
	imgPtr := flag.String("img", "", "original image path. (Required)")
	payloadPtr := flag.String("payload", "", "payload file. (Required)")
	outputImagePtr := flag.String("output", "modified_image.png", "output image path")
	flag.Parse()

	if !*decodePtr && !*insertPtr {
		fmt.Fprintf(os.Stderr, "You must specify an action: decode or insert\n")
		os.Exit(1)
	}

	if *decodePtr && *insertPtr {
		fmt.Fprintf(os.Stderr, "You must specify only one action: decode or insert\n")
		os.Exit(1)
	}

	if len(*imgPtr) == 0 {
		fmt.Fprintf(os.Stderr, "You must specify an image file\n")
		os.Exit(1)
	}

	myImage, err := pngint.ProcessImage(*imgPtr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed processing image: %s\n", err.Error())
		os.Exit(1)
	}

	if *decodePtr {
		parsedPayload, err := payload.ParsePayload(&myImage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to extract payload from image: %s\n", err.Error())
			os.Exit(1)
		}

		myPayload, err := payload.ExtractPayload(bytes.NewBuffer(parsedPayload))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to extract payload from image: %s\n", err.Error())
			os.Exit(1)
		}

		err = myPayload.WriteFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write payload file: %s\n", err.Error())
			os.Exit(1)
		}

	}

	if *insertPtr {
		f, err := os.Open(*payloadPtr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open payload file: %s\n", err.Error())
			os.Exit(1)
		}
		reader := bufio.NewReader(f)
		buf, err := ioutil.ReadAll(reader)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read payload file: %s\n", err.Error())
			os.Exit(1)
		}

		myPayload, err := payload.NewPayload(bytes.NewBuffer(buf), filepath.Ext(*payloadPtr)[1:])
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to generate payload data: %s\n", err.Error())
			os.Exit(1)
		}

		err = myPayload.InsertPayload(&myImage)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to insert payload into image: %s\n", err.Error())
			os.Exit(1)
		}
		err = myImage.UpdatePNG()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to insert payload into image: %s\n", err.Error())
			os.Exit(1)
		}
		myImage.ReconstructImage(*outputImagePtr)
	}

}
