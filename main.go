package main

import (
	"fmt"

	"github.com/ericoericochen/huffman/huffman"
)

func main() {
	text := "Helllo"
	// text := "HelloWorld"
	htree := huffman.CreateTreeFromText(text)

	fmt.Println(htree)
	// fmt.Println([]bool{true, false, true})
}