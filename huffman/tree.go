package huffman

import (
	"fmt"
	"slices"
)


type huffmanTree struct {
	root *huffmanNode
	bitmap map[string]Bits
}

type CodeNotFound struct {
	char string
}

func (e *CodeNotFound) Error() string {
	return fmt.Sprintf("CodeNotFound{char: %v}", e.char)
}

func (t *huffmanTree) getCode(char string) (Bits, error) {
	bits, ok := t.bitmap[char]
	if !ok {
		return nil, &CodeNotFound{char: char}
	}
	return bits, nil
}

func (t *huffmanTree) isRoot(n *huffmanNode) bool {
	return t.root == n
}

type huffmanNode struct {
	freq int
	char string
	left *huffmanNode
	right *huffmanNode
}

func (n huffmanNode) String() string {
	return fmt.Sprintf("Node{freq: %v, char: %v}", n.freq, n.char)
}

func (n *huffmanNode) isLeaf() bool {
	return n.left == nil && n.right == nil
}

func mergeNodes(a, b *huffmanNode) *huffmanNode {
	node := huffmanNode{}
	node.left = a
	node.right = b
	node.freq = a.freq + b.freq
	return &node
}

func preOrderTraversal(n *huffmanNode) {
	if n == nil {
		return
	}

	fmt.Println(n)
	preOrderTraversal(n.left)
	preOrderTraversal(n.right)
}

// get the mapping of char to their huffman codes (sequence of bits)
func getBitsForEachChar(n *huffmanNode) map[string]Bits  {
	bits := []Bit{}
	mapping := make(map[string]Bits)

	var traverse func(n *huffmanNode)
	traverse = func(n *huffmanNode) {
		if n.isLeaf() {
			mapping[n.char] = slices.Clone(bits)
			return
		}

		bits = append(bits, Zero)
		traverse(n.left)
		bits[len(bits) - 1] = One
		traverse(n.right)
		bits = bits[:len(bits) - 1]
	}

	traverse(n)
	return mapping
}


func createHuffmanTreeFromText(text string) *huffmanTree {
	fmt.Println(text)

	// count frequency of each char in text
	freqs := make(map[string]int)
	for _, r := range text {
		c := string(r)
		freqs[c] = freqs[c] + 1
	}

	// create leaf node for each char
	nodes := []*huffmanNode{}
	for char, freq := range freqs {
		node := huffmanNode{
			freq: freq,
			char: char,
		}
		nodes = append(nodes, &node)
	}

	// construct huffman tree by greedily merging the 2 lowest freq nodes
	// until there is one node
	for len(nodes) != 1 {
		slices.SortFunc(nodes, func(a, b *huffmanNode) int {
			return a.freq - b.freq
		})

		// create node that merges the 2 lowest freq nodes
		nodeA := nodes[0]
		nodeB := nodes[1]
		nodes = nodes[2:]

		merged := mergeNodes(nodeA, nodeB)
		nodes = append(nodes, merged)
	}

	root := nodes[0]
	// preOrderTraversal(root)

	// fmt.Println(nodes)

	bitmap := getBitsForEachChar(root)
	fmt.Println(bitmap)


	tree := huffmanTree{
		root: root,
		bitmap: bitmap,
	}

	return &tree
}