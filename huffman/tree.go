package huffman

import (
	"fmt"
	"slices"
)

type Bit bool
type Bits []Bit
const (
	Zero Bit = false
	One Bit = true
)

func (bs Bits) String() string {
	rep := ""
	for _, bit := range bs {
		if bit {
			rep += "1"
		} else {
			rep += "0"	
		}
	}
	return rep
}

type HuffmanTree struct {
	root *HuffmanNode
	mapping map[string]Bits
}

type HuffmanNode struct {
	count int
	char string
	left *HuffmanNode
	right *HuffmanNode
}

func (n HuffmanNode) String() string {
	return fmt.Sprintf("Node{count: %v, char: %v}", n.count, n.char)
}

func (n *HuffmanNode) IsLeaf() bool {
	return n.char != ""
}

func mergeNodes(a, b *HuffmanNode) *HuffmanNode {
	node := HuffmanNode{}
	node.left = a
	node.right = b
	node.count = a.count + b.count
	return &node
}

func preOrderTraversal(n *HuffmanNode) {
	if n == nil {
		return
	}

	fmt.Println(n)
	preOrderTraversal(n.left)
	preOrderTraversal(n.right)
}

// get the mapping of char to their huffman codes (sequence of bits)
func getBitsForEachChar(n *HuffmanNode) map[string]Bits  {
	bits := []Bit{}
	mapping := make(map[string]Bits)

	var traverse func(n *HuffmanNode)
	traverse = func(n *HuffmanNode) {
		if n.IsLeaf() {
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


func CreateTreeFromText(text string) *HuffmanTree {
	fmt.Println(text)

	// count frequency of each char in text
	counts := make(map[string]int)
	for _, r := range text {
		c := string(r)
		counts[c] = counts[c] + 1
	}

	// create leaf node for each char
	nodes := []*HuffmanNode{}
	for char, count := range counts {
		node := HuffmanNode{
			count: count,
			char: char,
		}
		nodes = append(nodes, &node)
	}

	// construct huffman tree by greedily merging the 2 lowest count nodes
	// until there is one node
	for len(nodes) != 1 {
		slices.SortFunc(nodes, func(a, b *HuffmanNode) int {
			return a.count - b.count
		})

		// create node that merges the 2 lowest count nodes
		nodeA := nodes[0]
		nodeB := nodes[1]
		nodes = nodes[2:]

		merged := mergeNodes(nodeA, nodeB)
		nodes = append(nodes, merged)
	}

	root := nodes[0]
	preOrderTraversal(root)

	// fmt.Println(nodes)

	mapping := getBitsForEachChar(root)
	fmt.Println(mapping)


	tree := HuffmanTree{}
	return &tree
}