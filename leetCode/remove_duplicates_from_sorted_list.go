// https://leetcode.com/problems/remove-duplicates-from-sorted-list/

package main

func deleteDuplicates(head *ListNode) *ListNode {
	ptr := head

	for ptr != nil && ptr.Next != nil {
		if ptr.Val == ptr.Next.Val {
			ptr.Next = ptr.Next.Next
		} else {
			ptr = ptr.Next
		}
	}

	return head
}
