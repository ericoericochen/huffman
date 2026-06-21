package huffman

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const FILE_EXTENSION = "hfcode"

type FileCompressor struct {}

func NewFileCompressor() *FileCompressor {
	return &FileCompressor{}
}

// errors
type InvalidHeaderError struct {}

func (e *InvalidHeaderError) Error() string {
	return fmt.Sprintf("Invalid %v header", FILE_EXTENSION)
}

// methods
func (fc *FileCompressor) Encode(fp string, saveFp string) {
	// open file
	file, err := os.Open(fp)
	
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// create huffman code
	huffmanCode := fc.CreateHuffmanCode(file)
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
		charBytes := make([]byte, 4)
		freqBytes := make([]byte, 8)
		binary.LittleEndian.PutUint32(charBytes, uint32(char))
		binary.LittleEndian.PutUint64(freqBytes, uint64(freq))
		writer.Write(charBytes)
		writer.Write(freqBytes)
	}

	// encode file contents using huffman tree
	var currentByte byte
	currentIdx := 0

	// stream encoded bits
	file.Seek(0, io.SeekStart)
	runesCh := streamRunes(file)
	bitsStream := huffmanCode.StreamEncode(runesCh)

	// set encoded bit in byte
	// when we reach end of byte, write to file and create new byte to write to
	for bit := range bitsStream {
		if bit {
			currentByte |= 1 << (7 - currentIdx)
		}
		currentIdx ++

		if currentIdx == 8 {
			writer.Write([]byte{currentByte})
			currentIdx = 0
			currentByte = 0
		}
	}

	// write byte we haven't written to disk yet
	// has additional padding
	if currentIdx > 0 {
		writer.Write([]byte{currentByte})
	}

	defer writer.Flush()
}


func (fc *FileCompressor) CreateHuffmanCode(f *os.File) *HuffmanCode {
	freqs := make(map[rune]int)
	f.Seek(0, io.SeekStart)
	for r := range streamRunes(f) {
		freqs[r]++
	}
	huffmanCode := NewHuffmanCodeFromFreq(freqs)
	return huffmanCode
}

func (fc *FileCompressor) Decode(fp string) {
	fmt.Println("decoding file")

	if !strings.HasSuffix(fp, FILE_EXTENSION) {
		log.Fatal("Not a .hfcode file")
	}
	
	file, err := os.Open(fp)
	if err != nil {
		log.Fatal(err)
	}

	huffmanCode, err := fc.DecodeHuffmanCodeFromHeader(file)
	if err != nil {
		log.Fatal(err)
	}

	saveFp := strings.TrimSuffix(fp, "." + FILE_EXTENSION)
	saveFile, err := os.Create(saveFp)
	if err != nil {
		log.Fatal(err)
	}
	writer := bufio.NewWriter(saveFile)
	defer writer.Flush()

	totalChars := huffmanCode.tree.TotalChars()
	currentChars := 0
	bitStream := streamBits(file)
	
	outputStream := huffmanCode.StreamDecode(bitStream)
	for r := range outputStream {
		char := string(r)
		currentChars ++
		writer.WriteString(char)
		
		if currentChars == totalChars {
			break
		}
	}
}


func (fc *FileCompressor) DecodeHuffmanCodeFromHeader(f *os.File) (*HuffmanCode, error) {
	var numCharsRaw uint32
	err := binary.Read(f, binary.LittleEndian, &numCharsRaw)
	if err != nil {
		return nil, &InvalidHeaderError{}
	}

	numChars := int(numCharsRaw)
	freqs := make(map[rune]int)

	for i := 0; i < numChars; i++ {
		var runeRaw uint32
		var freqRaw uint64
		
		err := binary.Read(f, binary.LittleEndian, &runeRaw)
		if err != nil {
			return nil, &InvalidHeaderError{}
		}

		err = binary.Read(f, binary.LittleEndian, &freqRaw)
		if err != nil {
			return nil, &InvalidDecodeError{}
		}

		char := rune(runeRaw)
		freq := int(freqRaw)
		freqs[char] = freq
	}

	huffmanCode := NewHuffmanCodeFromFreq(freqs)
	return huffmanCode, nil
}


func streamRunes(f *os.File) <-chan rune {
	ch := make(chan rune)
	reader := bufio.NewReader(f)

	go func() {
		defer close(ch)
		for {
			c, _, err := reader.ReadRune()
			if err != nil {
				break
			}
			ch <- c
		}
	}()

	return ch
}

func streamBits(f *os.File) <-chan Bit {
	reader := bufio.NewReader(f)
	ch := make(chan Bit)

	go func() {
		for {
			b, err := reader.ReadByte()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			bits := ByteToBits(b)
			for _, bit := range bits {
				ch <- bit
			}
		}
	}()

	return ch
}