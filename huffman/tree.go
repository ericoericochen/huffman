package huffman

import (
	"container/heap"
	"fmt"
	"slices"
)

type huffmanNode struct {
	freq int
	char rune
	left *huffmanNode
	right *huffmanNode
}

// helper functions
// check if 2 huffman trees are the same
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

	if node.IsLeaf() {
		ch <- node
		return
	}

	traverseLeafNodes(node.left, ch)
	traverseLeafNodes(node.right, ch)
}

func (n huffmanNode) String() string {
	return fmt.Sprintf("Node{freq: %v, char: %v}", n.freq, string(n.char))
}

func (n *huffmanNode) GetCharString() string {
	return string(n.char)
}

func (n *huffmanNode) IsLeaf() bool {
	return n.left == nil && n.right == nil
}

func (n *huffmanNode) Equals(other *huffmanNode) bool {
	return sameHuffmanTree(n, other)
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

func newHuffumanTree(charFreqs map[rune]int) *huffmanNode {
	// create leaf node for each char
	unicodes := make([]rune, 0, len(charFreqs))
	for k := range charFreqs {
		unicodes = append(unicodes, k)
	}

	slices.Sort(unicodes)
	
	// add to min priority queue where priority is frequency
	pq := NewPriorityQueue[*huffmanNode]()
	for _, unicode := range unicodes {
		freq := charFreqs[unicode]
		node := huffmanNode{
			freq: freq,
			char: unicode,
		}
		// nodes = append(nodes, &node)
		item := &PriorityQueueItem[*huffmanNode]{
			value: &node,
			priority: freq,
			// index: i,
		}
		heap.Push(&pq, item)
	}
	
	heap.Init(&pq)

	// construct huffman tree by greedily merging the 2 lowest freq nodes
	// until there is one node
	for pq.Len() > 1 {
		a := heap.Pop(&pq).(*PriorityQueueItem[*huffmanNode])
		b := heap.Pop(&pq).(*PriorityQueueItem[*huffmanNode])

		// create node that merges the 2 lowest freq nodes
		nodeA := a.value
		nodeB := b.value

		merged := mergeNodes(nodeA, nodeB)
		item := &PriorityQueueItem[*huffmanNode]{
			value: merged,
			priority: merged.freq,
		}
		heap.Push(&pq, item)
	}

	item := heap.Pop(&pq).(*PriorityQueueItem[*huffmanNode])
	root := item.value
	return root
}