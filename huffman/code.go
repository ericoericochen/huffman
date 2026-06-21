package huffman

import "fmt"


type HuffmanCode struct {
	tree *huffmanTree
}

type CompressionStats struct {
	OriginalBits int64
	CompressedBits int64
	CompressionRatio float32
}

// errors
type InvalidDecodeError struct {
	bits Bits
}

func (e *InvalidDecodeError) Error() string {
	return fmt.Sprintf("InvalidDecodeError - Unable to decode %v", e.bits)
}

// methods
func (hc *HuffmanCode) GetPrefixCode(char rune) (Bits, error) {
	return hc.tree.getCode(char)
}

// get the number of unique unicode characters
// in the huffman code
func (hc *HuffmanCode) NumChars() int {
	return len(hc.tree.bitmap)
}

func (hc *HuffmanCode) getTree() *huffmanTree {
	return hc.tree
}

func (hc *HuffmanCode) Equals(other *HuffmanCode) bool {
	return hc.tree.equals(other.tree)
}

func (hc *HuffmanCode) Encode(chars []rune) (Bits, error) {
	var bits Bits = []Bit{}

	// get code for each char and concat them together
	for _, r := range chars {
		code, err := hc.GetPrefixCode(r)
		if err != nil {
			return nil, err
		}

		bits = append(bits, code...)
	}
	
	return bits, nil
}

func (hc *HuffmanCode) StreamEncode(chars <-chan rune) <-chan Bit {
	ch := make(chan Bit)
	
	go func() {
		defer close(ch)
		for r := range chars {
			code, _ := hc.GetPrefixCode(r)
			for _, bit := range code {
				ch <- bit
			}
		}
	}()

	return ch
}

func (hc *HuffmanCode) Decode(bits Bits) (string, error) {
	decoded := ""
	node := hc.tree.root

	for _, bit := range bits {
		if bit == Zero {
			node = node.left
		} else if bit == One {
			node = node.right
		}

		// invalid decode if node is nil
		// went into part of tree where there is no prefix code
		// for the bits we're looking at
		if node == nil {
			return "", &InvalidDecodeError{bits: bits,}
		}

		if node.isLeaf() {
			decoded += node.GetCharString()
			node = hc.tree.root
		}
	}

	// successful decoding means last bit resulted in a decoded char
	// and node is reset to root node
	if !hc.tree.isRoot(node) {
		return "", &InvalidDecodeError{bits: bits,}
	}

	return decoded, nil
}

func (hc *HuffmanCode) StreamDecode(bits <-chan Bit) <-chan rune {
	ch := make(chan rune)

	go func() {
		defer close(ch)
		node := hc.tree.root
		for bit := range bits {
			if bit == Zero {
				node = node.left
			} else if bit == One {
				node = node.right
			}

			if node.isLeaf() {
				ch <- node.char
				node = hc.tree.root
			}
		}
	}()

	return ch
}

func (hc *HuffmanCode) GetCompressionStats() CompressionStats {
	originalBits := 0
	compressedBits := 0

	for node := range hc.tree.traverseLeafNodes() {
		charString := node.GetCharString()
		numBits := GetBits(charString) * node.freq
		code, _ := hc.GetPrefixCode(node.char)
		numCompressedBits := len(code) * node.freq
		originalBits += numBits
		compressedBits += numCompressedBits
	}

	// calculate compression ratio
	savedBits := originalBits - compressedBits
	compressionRatio := float32(savedBits) / float32(originalBits)

	stats := CompressionStats{
		OriginalBits: int64(originalBits),
		CompressedBits: int64(compressedBits),
		CompressionRatio: compressionRatio,
	}

	return stats
}

func (hc *HuffmanCode) GetEntropy() float32 {
	return 0
} 


func CreateHuffmanCodeFromText(text string) *HuffmanCode {
	tree := createHuffmanTreeFromText(text)
	hc := HuffmanCode{
		tree: tree,
	}

	return &hc
}

// freqs: frequency of unicode code point (rune)
func NewHuffmanCodeFromFreq(freqs map[rune]int) *HuffmanCode {
	tree := newHuffumanTreeFromFreq(freqs)
	hc := HuffmanCode{tree: tree,}
	return &hc
}