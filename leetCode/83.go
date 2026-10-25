// https://leetcode.com/problems/remove-duplicates-from-sorted-list/
// N is the number of nodes
// time complexity is O(N)
// space complexity is O(1)

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
