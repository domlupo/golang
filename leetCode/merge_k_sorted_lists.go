// https://leetcode.com/problems/merge-k-sorted-lists
// N is the total number of nodes, K is the number of linked lists
// time complexity is O(N*log(K))
// space complexity is O(K)

package main

import "container/heap"

func mergeKLists(lists []*ListNode) *ListNode {
	h := &Heap{}
	heap.Init(h)
	for _, list := range lists {
		curr := list
		if curr != nil {
			heap.Push(h, curr)
		}
	}

	dummy := &ListNode{}
	curr := dummy
	for h.Len() > 0 {
		n := heap.Pop(h).(*ListNode)
		curr.Next = n
		curr = curr.Next
		if n.Next != nil {
			heap.Push(h, n.Next)
		}
	}

	return dummy.Next
}

type Heap []*ListNode

func (h Heap) Len() int           { return len(h) }
func (h Heap) Less(i, j int) bool { return h[i].Val < h[j].Val }
func (h Heap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *Heap) Push(x any) {
	*h = append(*h, x.(*ListNode))
}
func (h *Heap) Pop() any {
	old := *h
	n := len(old) - 1
	x := old[n]
	*h = old[:n]
	return x
}
