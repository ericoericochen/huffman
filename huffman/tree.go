package huffman

import (
	"fmt"
	"slices"
)


type huffmanTree struct {
	root *huffmanNode
	bitmap map[rune]Bits
	leafNodes map[rune]*huffmanNode
}

type CodeNotFound struct {
	char rune
}

func (e *CodeNotFound) Error() string {
	return fmt.Sprintf("CodeNotFound{char: %v}", string(e.char))
}

func (t *huffmanTree) getCode(char rune) (Bits, error) {
	bits, ok := t.bitmap[char]
	if !ok {
		return nil, &CodeNotFound{char: char}
	}
	return bits, nil
}

func (t *huffmanTree) isRoot(n *huffmanNode) bool {
	return t.root == n
}

func (t *huffmanTree) getLeafNode(char rune) *huffmanNode {
	return t.leafNodes[char]
}

func (t *huffmanTree) traverseLeafNodes() <-chan *huffmanNode  {
	ch := make(chan *huffmanNode)

	go func() {
		traverseLeafNodes(t.root, ch)
		close(ch)
	}()

	return ch
}

func (t *huffmanTree) equals(other *huffmanTree) bool {
	return sameHuffmanTree(t.root, other.root)
}

func sameHuffmanTree(a, b *huffmanNode) bool {
	// check both nodes are nil
	if a == nil || b == nil {
		return a == b
	}

	// check current node
	if !a.equals(b) {
		return false
	}

	// check left subtree
	if !sameHuffmanTree(a.left, b.left) {
		return false
	}

	// check right subtree
	if !sameHuffmanTree(a.right, b.right) {
		return false
	}

	return true
}

func traverseLeafNodes(node *huffmanNode, ch chan<- *huffmanNode) {
	if node == nil {
		return
	}

	if node.isLeaf() {
		ch <- node
		return
	}

	traverseLeafNodes(node.left, ch)
	traverseLeafNodes(node.right, ch)
}

type huffmanNode struct {
	freq int
	char rune
	left *huffmanNode
	right *huffmanNode
}

func (n huffmanNode) String() string {
	return fmt.Sprintf("Node{freq: %v, char: %v}", n.freq, string(n.char))
}

func (n *huffmanNode) GetCharString() string {
	return string(n.char)
}

func (n *huffmanNode) isLeaf() bool {
	return n.left == nil && n.right == nil
}

func (n *huffmanNode) equals(other *huffmanNode) bool {
	if n.freq != other.freq {
		return false
	}

	if n.char != other.char {
		return false
	}
	
	return true
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

func traverse(n *huffmanNode, f func(*huffmanNode)) {
	if n == nil {
		return
	}

	f(n)
	traverse(n.left, f)
	traverse(n.right, f)
}

// get the mapping of char to their huffman codes (sequence of bits)
func getBitsForEachChar(n *huffmanNode) map[rune]Bits  {
	bits := []Bit{}
	mapping := make(map[rune]Bits)

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
	// count frequency of each char in text
	freqs := make(map[rune]int)
	for _, r := range text {
		freqs[r] = freqs[r] + 1
	}

	return newHuffumanTreeFromFreq(freqs)
}

func newHuffumanTreeFromFreq(freqs map[rune]int) *huffmanTree {
	// create leaf node for each char
	nodes := []*huffmanNode{}
	unicodes := make([]rune, 0, len(freqs))
	for k := range freqs {
		unicodes = append(unicodes, k)
	}

	slices.Sort(unicodes)
	for _, unicode := range unicodes {
		freq := freqs[unicode]
		node := huffmanNode{
			freq: freq,
			char: unicode,
		}
		nodes = append(nodes, &node)
	}

	// construct huffman tree by greedily merging the 2 lowest freq nodes
	// until there is one node
	for len(nodes) != 1 {
		// TODO: replace this with min priority queue
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
	bitmap := getBitsForEachChar(root)
	leafNodes := make(map[rune]*huffmanNode)

	traverse(root, func(n *huffmanNode) {
		if n.isLeaf() {
			leafNodes[n.char] = n
		}
	})

	// fmt.Println(bitmap)

	tree := huffmanTree{
		root: root,
		bitmap: bitmap,
		leafNodes: leafNodes,
	}

	return &tree
}