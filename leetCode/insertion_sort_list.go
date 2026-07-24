// https://leetcode.com/problems/insertion-sort-list/
// N is the numbers of nodes
// time complexity is O(N * N)
// space complexity is O(1)

package main

func insertionSortList(head *ListNode) *ListNode {
	dummy := &ListNode{}
	curr := head
	for curr != nil {
		prev := dummy
		for prev.Next != nil && prev.Next.Val <= curr.Val {
			prev = prev.Next
		}
		next := curr.Next
		curr.Next = prev.Next
		prev.Next = curr
		curr = next
	}
	return dummy.Next
}
