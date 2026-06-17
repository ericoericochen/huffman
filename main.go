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

	fmt.Println("same: ", text == decoded)

	// compression stats
	stats := hc.GetCompressionStats()
	fmt.Println(stats)

	// measure bits saved
	fmt.Println("original bits: ", stats.OriginalBits)
	fmt.Println("compressed bits: ", stats.CompressedBits)
	fmt.Println("compression ratio: ", stats.CompressionRatio)
}