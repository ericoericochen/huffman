package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type HfcodeCli struct {}

func streamRunes(f *os.File) <-chan rune {
	ch := make(chan rune)
	f.Seek(0, io.SeekStart)
	scanner := bufio.NewScanner(f)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			for _, r := range line {
				ch <- r
			}
		}
		close(ch)
	}()

	return ch
}

// encode a file into a .hfcode file
func (cli *HfcodeCli) Encode(fp string) {
	fmt.Println("encoding ", fp)

	file, err := os.Open(fp)
	
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// frequency of unicode characters
	freqs := make(map[rune]int)
	for r := range streamRunes(file) {
		freqs[r]++
	}

	fmt.Println(freqs)
	fmt.Println(len(freqs))
}

// decode a hfcode file
func (cli *HfcodeCli) Decode(fp string) {

}


func main() {
	fmt.Println("hi from cli")

	cli := HfcodeCli{}

	args := os.Args
	cmd := args[1]

	switch cmd {
	case "encode":
		fmt.Println("encode")
		cli.Encode(args[2])
	default:
		fmt.Println("non command found")
	}

	// const s = "สวัสดี"


	// fmt.Println(len(s))
	//     fmt.Println("Rune count:", utf8.RuneCountInString(s))
}