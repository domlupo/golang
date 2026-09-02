// https://leetcode.com/problems/last-stone-weight/
// N is the length of stones
// time complexity is O(N*log(N))
// space complexity is O(N)

package main

import "container/heap"

func lastStoneWeight(stones []int) int {
	h := &MaxHeap1046{}
	heap.Init(h)
	for _, s := range stones {
		heap.Push(h, s)
	}

	for h.Len() >= 2 {
		s1 := heap.Pop(h).(int)
		s2 := heap.Pop(h).(int)
		if s1 != s2 {
			heap.Push(h, s1-s2)
		}
	}

	if h.Len() == 1 {
		return (*h)[0]
	} else {
		return 0
	}
}

type MaxHeap1046 []int

func (h MaxHeap1046) Len() int           { return len(h) }
func (h MaxHeap1046) Less(i, j int) bool { return h[i] > h[j] }
func (h MaxHeap1046) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap1046) Push(x any) {
	*h = append(*h, x.(int))
}
func (h *MaxHeap1046) Pop() any {
	old := *h
	n := len(old) - 1
	x := old[n]
	*h = old[:n]
	return x
}
