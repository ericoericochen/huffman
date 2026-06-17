package huffman

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

type FileCompressor struct {}

func NewFileCompressor() *FileCompressor {
	return &FileCompressor{}
}

func (fc *FileCompressor) Encode(fp string, saveFp string) {
	// open file
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

	// construct huffman code
	huffmanCode := NewHuffmanCodeFromFreq(freqs)
	tree := huffmanCode.getTree()

	// create save file
	saveFile, err := os.Create(saveFp)
	if err != nil {
		log.Fatal(err)
	}

	defer saveFile.Close()

	// construct header
	// save the number of unique chars in uint32
	writer := bufio.NewWriter(saveFile)
	numCharsBs := make([]byte, 4)
	numChars := huffmanCode.NumChars()
	binary.LittleEndian.PutUint32(numCharsBs, uint32(numChars))	
	writer.Write(numCharsBs)

	// for each char, save its unicode point and frequency
	for node := range tree.traverseLeafNodes() {
		freq := node.freq
		char := node.char
		fmt.Println(string(char), ": ", freq)
		charBytes := make([]byte, 4)
		freqBytes := make([]byte, 8)
		binary.LittleEndian.PutUint32(charBytes, uint32(char))
		binary.LittleEndian.PutUint64(freqBytes, uint64(freq))
		writer.Write(charBytes)
		writer.Write(freqBytes)
	}

	defer writer.Flush()
}

func (fc *FileCompressor) Decode(fp string) {

}



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