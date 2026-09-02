// https://leetcode.com/problems/kth-largest-element-in-an-array/
// N is the length of nums, K is passed k
// time complexity is O(N * log(K))
// space complexity is O(K)

package main

import "container/heap"

func findKthLargest(nums []int, k int) int {
	h := &MinHeap703{}
	for _, n := range nums {
		heap.Push(h, n)
		if h.Len() > k {
			heap.Pop(h)
		}
	}
	return heap.Pop(h).(int)
}
