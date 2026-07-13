// https://leetcode.com/problems/remove-duplicates-from-sorted-list-ii/
// N is the length of linked list
// time complexity is O(N)
// space complexity is O(1)

package main

func deleteDuplicates82(head *ListNode) *ListNode {
	dummy := &ListNode{Next: head, Val: -101}
	prev := dummy
	curr := prev.Next
	for curr != nil && curr.Next != nil {
		if curr.Val != curr.Next.Val {
			prev = curr
			curr = curr.Next
			continue
		}

		for curr.Next != nil && curr.Val == curr.Next.Val {
			curr = curr.Next
		}

		prev.Next = curr.Next
		curr = prev
	}

	return dummy.Next
}
