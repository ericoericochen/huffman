package huffman

import "fmt"


type HuffmanCode struct {
	tree *huffmanNode
	bitmap map[rune]Bits
	leafNodes map[rune]*huffmanNode
	charFreqs map[rune]int
}

type CompressionStats struct {
	OriginalBits int64
	CompressedBits int64
	CompressionRatio float32
	CompressionPercentage float32
	BitsPerSymbol float32
}

type StreamEncodeResult struct {
	Val Bit
	Err error
}

type StreamDecodeResult struct {
	Val rune
	Err error
}

// errors
type InvalidDecodeError struct {
	bits Bits
}

func (e *InvalidDecodeError) Error() string {
	return fmt.Sprintf("InvalidDecodeError - Unable to decode %v", e.bits)
}

type CodeNotFound struct {
	char rune
}

func (e *CodeNotFound) Error() string {
	return fmt.Sprintf("CodeNotFound{char: %v}", string(e.char))
}


// methods
func (hc *HuffmanCode) GetPrefixCode(char rune) (Bits, error) {
	bits, ok := hc.bitmap[char]
	if !ok {
		return nil, &CodeNotFound{char: char}
	}
	return bits, nil
}

// get the number of unique unicode characters
// in the huffman code
func (hc *HuffmanCode) NumChars() int {
	return len(hc.bitmap)
}

// get the total number of characters
func (hc *HuffmanCode) TotalChars() int {
	totalChars := 0
	for leafNode := range hc.LeafNodes() {
		totalChars += leafNode.freq
	}
	return totalChars
}

// get the huffman tree
func (hc *HuffmanCode) Tree() *huffmanNode {
	return hc.tree
}

// get the leaf nodes
func (hc *HuffmanCode) LeafNodes() <-chan *huffmanNode {
	ch := make(chan *huffmanNode)
	go func() {
		defer close(ch)
		for _, ln := range hc.leafNodes {
			ch <- ln
		}
	}()
	return ch
}

func (hc *HuffmanCode) Equals(other *HuffmanCode) bool {
	return hc.tree.Equals(other.tree)
}

// get bit map
func (hc *HuffmanCode) BitMap() map[rune]Bits {
	return hc.bitmap
}

func (hc *HuffmanCode) Encode(chars []rune) (Bits, error) {
	var bits Bits = []Bit{}

	charsCh := make(chan rune)
	go func() {
		defer close(charsCh)
		for _, r := range chars {
			charsCh <- r
		}
	}()

	for result := range hc.StreamEncode(charsCh) {
		if result.Err != nil {
			return []Bit{}, result.Err
		}
		bits = append(bits, result.Val)
	}
	return bits, nil
}

func (hc *HuffmanCode) EncodeText(text string) (Bits, error) {
	chars := []rune(text)
	return hc.Encode(chars)
} 

func (hc *HuffmanCode) StreamEncode(charsCh <-chan rune) <-chan *StreamEncodeResult {
	ch := make(chan *StreamEncodeResult)

	go func() {
		defer close(ch)
		for r := range charsCh {
			code, err := hc.GetPrefixCode(r)
			if err != nil {
				ch <- &StreamEncodeResult{Err: err}
				break
			}
			for _, bit := range code {
				ch <- &StreamEncodeResult{Val: bit}
			}
		}
	}()
	return ch
}

func (hc *HuffmanCode) Decode(bits Bits) ([]rune, error) {
	decoded := []rune{}
	bitsCh := make(chan Bit)

	go func() {
		defer close(bitsCh)
		for _, bit := range bits {
			bitsCh <- bit
		}
	}()

	for result := range hc.StreamDecode(bitsCh) {
		if result.Err != nil {
			return []rune{}, result.Err
		}
		decoded = append(decoded, result.Val)
	}

	return decoded, nil
}

func (hc *HuffmanCode) DecodeText(bits Bits) (string, error) {
	output, err := hc.Decode(bits)
	text := string(output)
	return text, err
}

func (hc *HuffmanCode) StreamDecode(bitsCh <-chan Bit) <-chan *StreamDecodeResult {
	ch := make(chan *StreamDecodeResult)

	go func() {
		defer close(ch)
		node := hc.tree
		for bit := range bitsCh {
			if bit == Zero {
				node = node.left
			} else if bit == One {
				node = node.right
			}

			if node == nil {
				ch <- &StreamDecodeResult{Err: &InvalidDecodeError{}}
				break
			}

			if node.IsLeaf() {
				ch <- &StreamDecodeResult{Val: node.char}
				node = hc.tree
			}
		}
	}()

	return ch
}

func (hc *HuffmanCode) GetCompressionStats() CompressionStats {
	originalBits := 0
	compressedBits := 0
	numSymbols := 0

	for node := range hc.LeafNodes() {
		charString := node.GetCharString()
		numBits := GetBits(charString) * node.freq
		code, _ := hc.GetPrefixCode(node.char)
		numCompressedBits := len(code) * node.freq
		originalBits += numBits
		compressedBits += numCompressedBits
		numSymbols += node.freq
	}

	// calculate stats
	compressionRatio := float32(originalBits) / float32(compressedBits)
	compressionPercentage := 1 - float32(compressedBits) / float32(originalBits)
	bitsPerSymbol := float32(compressedBits) / float32(numSymbols)

	stats := CompressionStats{
		OriginalBits: int64(originalBits),
		CompressedBits: int64(compressedBits),
		CompressionRatio: compressionRatio,
		CompressionPercentage: compressionPercentage,
		BitsPerSymbol: bitsPerSymbol,

	}

	return stats
}

func (hc *HuffmanCode) GetEntropy() float32 {
	return 0
} 


func countCharFreqs(text string) map[rune]int {
	charFreqs := make(map[rune]int)
	for _, r := range text {
		charFreqs[r] = charFreqs[r] + 1
	}
	return charFreqs
}

func NewHuffmanCodeFromText(text string) *HuffmanCode {
	charFreqs := countCharFreqs(text)
	return NewHuffmanCode(charFreqs)
}


// creates a huffman code from character frequencies
// charFreqs - frequency of unicode codepoints
func NewHuffmanCode(charFreqs map[rune]int) *HuffmanCode {
	tree := newHuffumanTree(charFreqs)

	// map of unicode codepoint to bits encoding
	bitmap := getBitsForEachChar(tree)

	// map unicode codepoint to leaf node in huffman tree
	leafNodes := make(map[rune]*huffmanNode)
	traverse(tree, func(n *huffmanNode) {
		if n.IsLeaf() {
			leafNodes[n.char] = n
		}
	})
	
	hc := HuffmanCode{
		tree: tree,
		leafNodes: leafNodes,
		bitmap: bitmap,
		charFreqs: charFreqs,
		
	}
	return &hc
}