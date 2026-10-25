// https://leetcode.com/problems/merge-two-sorted-lists/
// M is number of list1 nodes, N is the numbers of list2 nodes.
// time complexity is O(M + N)
// space complexity is O(1)

package main

type ListNode struct {
	Val  int
	Next *ListNode
}

func mergeTwoSortedLists(list1 *ListNode, list2 *ListNode) *ListNode {
	dummy := &ListNode{}
	curr := dummy
	p1, p2 := list1, list2

	for p1 != nil && p2 != nil {
		if p1.Val <= p2.Val {
			curr.Next = p1
			p1 = p1.Next
		} else {
			curr.Next = p2
			p2 = p2.Next
		}
		curr = curr.Next
	}
	if p1 != nil {
		curr.Next = p1
	}
	if p2 != nil {
		curr.Next = p2
	}

	return dummy.Next
}
