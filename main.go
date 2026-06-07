package main

import "fmt"

func main() {
	text := "HelloWorld"
	htree := CreateHuffmanTreeFromText(text)

	fmt.Println(htree)
}