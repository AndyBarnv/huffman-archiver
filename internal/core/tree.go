package core

import "container/heap"

type Node struct {
	Char  byte
	Freq  uint64
	Left  *Node
	Right *Node
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Freq < pq[j].Freq
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*Node))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func BuildTree(freq map[byte]uint64) *Node {
	pq := make(PriorityQueue, 0, len(freq))
	heap.Init(&pq)

	for char, count := range freq {
		heap.Push(&pq, &Node{Char: char, Freq: count})
	}

	if pq.Len() == 0 {
		return nil
	}

	if pq.Len() == 1 {
		node := heap.Pop(&pq).(*Node)
		return &Node{Freq: node.Freq, Left: node}
	}

	for pq.Len() > 1 {
		left := heap.Pop(&pq).(*Node)
		right := heap.Pop(&pq).(*Node)
		parent := &Node{
			Freq:  left.Freq + right.Freq,
			Left:  left,
			Right: right,
		}
		heap.Push(&pq, parent)
	}

	return heap.Pop(&pq).(*Node)
}
