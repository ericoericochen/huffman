package main

import (
	"fmt"
	"os"

	"github.com/ericoericochen/huffman/huffman"
)

type HfcodeCli struct {}

// encode a file into a .hfcode file
func (cli *HfcodeCli) Encode(fp string, saveFp string) {
	fileCompressor := huffman.NewFileCompressor()
	fileCompressor.Encode(fp, saveFp)	
}

// decode a hfcode file
func (cli *HfcodeCli) Decode(fp string) {
	fileCompressor := huffman.NewFileCompressor()
	fileCompressor.Decode(fp)
}


func main() {
	fmt.Println("hi from cli")

	cli := HfcodeCli{}

	args := os.Args
	cmd := args[1]

	switch cmd {
	case "encode":
		fmt.Println("encode")
		cli.Encode(args[2], args[3])
	case "decode": 
		cli.Decode(args[2])
	case "stats":
		fmt.Println("stats")
	default:
		fmt.Println("non command found")
	}
}