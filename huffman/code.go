package huffman

import "fmt"

const BITS_IN_BYTE = 8

type HuffmanCode struct {
	tree *huffmanTree
}

// errors
type InvalidDecodeError struct {
	bits Bits
}

func (e *InvalidDecodeError) Error() string {
	return fmt.Sprintf("InvalidDecodeError - Unable to decode %v", e.bits)
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

func MeasureCompressionRatio(text string, bits Bits) float32 {
	numOriginalBits := len(text) * BITS_IN_BYTE
	numCompressedBits := len(bits)
	numSavedBits := numOriginalBits - numCompressedBits
	compressionRatio := float32(numSavedBits) / float32(numOriginalBits)
	return compressionRatio
}

func CreateHuffmanCodeFromText(text string) *HuffmanCode {
	tree := createHuffmanTreeFromText(text)
	hc := HuffmanCode{
		tree: tree,
	}

	return &hc
}