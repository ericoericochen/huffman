package main

import (
	"fmt"

	"github.com/ericoericochen/huffman/huffman"
)

func main() {
	// text := "Helllo"
	text := "Hellllllllllllllllo"

	hc := huffman.CreateHuffmanCodeFromText(text)
	fmt.Println(hc)

	encoded, _ := hc.Encode(text)
	fmt.Println(encoded)

	decoded, _ := hc.Decode(encoded)
	fmt.Println(decoded)

	// measure bits saved
	fmt.Println("original bits: ", len(text) * huffman.BITS_IN_BYTE)
	fmt.Println("compressed bits: ", len(encoded))
	fmt.Println("compression ratio: ", huffman.MeasureCompressionRatio(text, encoded))
}