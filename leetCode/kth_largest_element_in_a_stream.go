// https://leetcode.com/problems/kth-largest-element-in-a-stream/

package main

import "container/heap"

type KthLargest struct {
	k       int
	minHeap *IntHeap
}

func Constructor3(k int, nums []int) KthLargest {
	minHeap := &IntHeap{}
	heap.Init(minHeap)
	kthLargest := KthLargest{k: k, minHeap: minHeap}
	for _, num := range nums {
		kthLargest.Add(num)
	}
	return kthLargest
}

func (this *KthLargest) Add(val int) int {
	if this.minHeap.Len() < this.k {
		heap.Push(this.minHeap, val)
	} else if val > (*this.minHeap)[0] {
		heap.Pop(this.minHeap)
		heap.Push(this.minHeap, val)
	}
	return (*this.minHeap)[0]
}

// IntHeap is a min-heap of ints.
type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
