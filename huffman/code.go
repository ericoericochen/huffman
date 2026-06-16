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
func (hc *HuffmanCode) GetPrefixCode(char string) (Bits, error) {
	return hc.tree.getCode(char)
}

func (hc *HuffmanCode) Encode(text string) (Bits, error) {
	var bits Bits = []Bit{}

	// get code for each char and concat them together
	for _, r := range text {
		char := string(r)
		code, err := hc.tree.getCode(char)
		if err != nil {
			return nil, err
		}

		bits = append(bits, code...)
	}
	
	return bits, nil
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
			decoded += node.char
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

func (hc *HuffmanCode) GetCompressionStats() CompressionStats {
	originalBits := 0
	compressedBits := 0

	for node := range hc.tree.traverseLeafNodes() {
		char := node.char
		numBits := GetBits(char) * node.freq
		code, _ := hc.GetPrefixCode(char)
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

	hc := HuffmanCode{}
	return &hc
}