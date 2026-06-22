package huffman

type PriorityQueue[T any] struct {

}

func (pq *PriorityQueue[T]) Push(item T, priority int) {}

func (pq *PriorityQueue[T]) Pop() (T, int) {}

func (pq *PriorityQueue[T]) Len() int {}

