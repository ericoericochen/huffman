package main

import (
	"fmt"
	"slices"
)

type huffmanTree struct {
	root *huffmanNode
	// mapping map[string]
}

type huffmanNode struct {
	count int
	char *string
	left *huffmanNode
	right *huffmanNode
}

func (n *huffmanNode) IsLeaf() bool {
	return n.char != nil
}

// func mergeHuffmanNodes(a, b *huffmanNode) *huffmanNode {

// }


func CreateHuffmanTreeFromText(text string) *huffmanTree {
	fmt.Println(text)

	// count frequency of each char in text
	counts := make(map[string]int)
	for _, r := range text {
		c := string(r)
		counts[c] = counts[c] + 1
	}

	// create leaf node for each char
	nodes := []*huffmanNode{}
	for char, count := range counts {
		node := huffmanNode{
			count: count,
			char: &char,
		}
		nodes = append(nodes, &node)
	}

	// construct huffman tree by greedily merging the 2 lowest count nodes
	// until there is one node
	for len(nodes) != 1 {
		fmt.Println("want to merge haha")
		slices.SortFunc(nodes, func(a, b *huffmanNode) int {
			return a.count - b.count
		})
		
		for _, n := range nodes {
			fmt.Println("count: ", n.count)
		} 
		break
	}


	tree := huffmanTree{}
	return &tree
}