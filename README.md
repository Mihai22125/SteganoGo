# SteganoGo


## Overview
`steganoGo` is a simple CLI capable of hiding any file within an PNG image.
Used technique is known as LSB (Least Significant Bit) [steganography](https://en.wikipedia.org/wiki/steganography) 

## Demonstration

| Original image                         | Payload                             | Result                                               |
| ---------------------------------------| ------------------------------------|------------------------------------------------------|
| ![Original File](examples/grass2.png)  | ![Payload file](examples/city.jpg)  | ![Encoded File](examples/modified_image.png)         |


The `Result` file contains the `Payload` file hidden in it. And as you can see it is fully transparent.

## Instalation

Clone the code.

```bash
git clone https://github.com/Mihai22125/SteganoGo
```


Get into the source directory.

```bash
cd SteganoGo/source
```

Build the code from the root directory

```bash
go build -o steganoGo
```

> This sends the output of `go build` to a file called `steganoGo` in the same directory.

## Usage

#### Encoding
```
steganoGo -insert -img <file-name> -payload <file-name> -output <file-name>
```
When encoding, the file with name given to flag `-payload` is hidden inside the file with the name given to flag `-img` and the resulting file is saved in a new file under name given to flag `-output`

#### Decoding
```
steganoGo -decode -img <file-name> 
```
When decoding, given file name of an image with previously encoded data in it, the data is extracted and saved in a new file in the current working directory.
The result file will have the same extension as when it was encoded.

## Disclaimer

The only supported format is PNG image format.


