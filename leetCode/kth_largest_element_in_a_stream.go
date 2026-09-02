// https://leetcode.com/problems/kth-largest-element-in-a-stream/
// K is passed k, N is number of times Add is called.
// time complexity for Add called N times is O(N*log(K)).
// space complexity is O(K)

package main

import "container/heap"

type KthLargest struct {
	HeapSize   int
	MinHeap703 *MinHeap703
}

func Constructor703(k int, nums []int) KthLargest {
	h := &MinHeap703{}
	heap.Init(h)
	for _, n := range nums {
		heap.Push(h, n)
		if h.Len() > k {
			heap.Pop(h)
		}
	}
	return KthLargest{
		HeapSize:   k,
		MinHeap703: h,
	}
}

func (this *KthLargest) Add(val int) int {
	heap.Push(this.MinHeap703, val)
	if this.MinHeap703.Len() > this.HeapSize {
		heap.Pop(this.MinHeap703)
	}
	return (*this.MinHeap703)[0]
}

type MinHeap703 []int

func (h MinHeap703) Len() int           { return len(h) }
func (h MinHeap703) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap703) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap703) Push(x any) {
	*h = append(*h, x.(int))
}
func (h *MinHeap703) Pop() any {
	old := *h
	n := len(old) - 1
	x := old[n]
	*h = old[:n]
	return x
}
