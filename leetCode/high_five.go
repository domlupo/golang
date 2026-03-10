// https://leetcode.com/problems/high-five/

package main

import "container/heap"

type IntHeap2 []int

// max heap, for min heap swap equality
func (h IntHeap2) Less(i, j int) bool { return h[i] > h[j] }
func (h IntHeap2) Len() int           { return len(h) }
func (h IntHeap2) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap2) Push(x any) {
	*h = append(*h, x.(int))
}
func (h *IntHeap2) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func highFive(items [][]int) [][]int {
	x := 5
	allScores := make(map[int]*IntHeap2)

	for _, item := range items {
		id := item[0]
		s := item[1]
		scores, ok := allScores[id]
		if !ok {
			scores = &IntHeap2{}
		}
		heap.Push(scores, s)
		allScores[id] = scores
	}

	var ans [][]int
	for i := 1; i <= 1000; i++ {
		s, ok := allScores[i]
		var sum int
		if ok {
			for range x {
				val := heap.Pop(s).(int)
				sum += val
			}
			ans = append(ans, []int{i, sum / x})
		}
	}
	return ans
}
