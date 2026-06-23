package main

import (
	"fmt"

	"github.com/ericoericochen/huffman/huffman"
)

func main() {
	// text := "Helllo"
	text := "Helllllllllllllllllllllllo"

	hc := huffman.NewHuffmanCodeFromText(text)
	fmt.Println(hc.BitMap())

	for r, b := range hc.BitMap() {
		fmt.Println(string(r), b)
	}

	encoded, _ := hc.EncodeText(text)
	fmt.Println(encoded)

	decoded, _ := hc.DecodeText(encoded)
	fmt.Println(decoded)

	fmt.Println("same: ", text == decoded)

	// compression stats
	stats := hc.GetCompressionStats()
	fmt.Println(stats)

	// charsCh := make(chan rune)

	// go func() {
	// 	for _, r := range chars {
	// 		charsCh <- r
	// 	}
	// 	close(charsCh)
	// }()

	// for bit := range hc.StreamEncode(charsCh) {
	// 	fmt.Println(bit)
	// }

	// // measure bits saved
	// fmt.Println("original bits: ", stats.OriginalBits)
	// fmt.Println("compressed bits: ", stats.CompressedBits)
	// fmt.Println("compression ratio: ", stats.CompressionRatio)

	// fc := huffman.NewFileCompressor()
	// fp := "./examples/tiny-shakespeare.txt"
	// saveFp := "./archive.hfcode"

	// file, _ := os.Open(fp)
	// saveFile, _ := os.Open(saveFp)

	// defer file.Close()
	// defer saveFile.Close()

	// hc1 := fc.CreateHuffmanCode(file)
	// fmt.Println("num chars hc1: ", hc1.NumChars())
	// hc2, _ := fc.DecodeHuffmanCodeFromHeader(saveFile)
	// equalHc := hc1.Equals(hc2)

	// fmt.Println("equal hc: ")
	// fmt.Println(equalHc)
}