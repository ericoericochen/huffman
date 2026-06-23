package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ericoericochen/huffman/huffman"
)

type HfcodeCli struct{}

// encode a file into a .hfcode file
func (cli *HfcodeCli) Encode(fp string, saveFp string) {
	fileCompressor := huffman.NewFileCompressor()
	fileCompressor.Encode(fp, saveFp)
}

// decode a hfcode file
func (cli *HfcodeCli) Decode(fp string, outputFp string) {
	fileCompressor := huffman.NewFileCompressor()
	fileCompressor.Decode(fp, outputFp)
}

// print stats about a .hfcode file
func (cli *HfcodeCli) Stats(fp string) {
	fileCompressor := huffman.NewFileCompressor()

	file, err := os.Open(fp)
	if err != nil {
		log.Fatal(err)
	}
	huffmanCode, err := fileCompressor.DecodeHuffmanCodeFromHeader(file)
	if err != nil {
		log.Fatal(err)
	}
	stats := huffmanCode.GetCompressionStats()
	printCompressionStats(stats)
}

func printCompressionStats(stats huffman.CompressionStats) {
	fmt.Printf("Original bits: %d bits\n", stats.OriginalBits)
	fmt.Printf("Compressed bits: %d bits\n", stats.CompressedBits)
	fmt.Printf("Compression ratio: %.2f\n", stats.CompressionRatio)
	fmt.Printf("Compression percentage: %.2f%%\n", stats.CompressionPercentage*100)
	fmt.Printf("Bits per symbol: %.2f bits/symbol\n", stats.BitsPerSymbol)
}

func main() {
	cli := HfcodeCli{}

	args := os.Args
	cmd := args[1]

	switch cmd {
	case "encode":
		cli.Encode(args[2], args[3])
	case "decode":
		cli.Decode(args[2], args[3])
	case "stats":
		cli.Stats(args[2])
	default:
		fmt.Println("non command found")
	}
}
