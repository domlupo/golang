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
	ptr := dummy

	for list1 != nil && list2 != nil {
		if list1.Val < list2.Val {
			ptr.Next = list1
			list1 = list1.Next
		} else {
			ptr.Next = list2
			list2 = list2.Next
		}
		ptr = ptr.Next
	}

	if list1 != nil {
		ptr.Next = list1
	}
	if list2 != nil {
		ptr.Next = list2
	}

	return dummy.Next
}
