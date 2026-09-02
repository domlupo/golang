// https://leetcode.com/problems/kth-largest-element-in-an-array/
// N is the length of nums, K is passed k
// time complexity is O(N * log(K))
// space complexity is O(K)

package main

import "container/heap"

func findKthLargest(nums []int, k int) int {
	h := &MinHeap215{}
	heap.Init(h)
	for _, n := range nums {
		heap.Push(h, n)
		if h.Len() > k {
			heap.Pop(h)
		}
	}
	return (*h)[0]
}

type MinHeap215 []int

func (h MinHeap215) Len() int           { return len(h) }
func (h MinHeap215) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap215) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap215) Push(x any)        { *h = append(*h, x.(int)) }
func (h *MinHeap215) Pop() any {
	old := *h
	n := len(old) - 1
	x := old[n]
	*h = old[:n]
	return x
}
