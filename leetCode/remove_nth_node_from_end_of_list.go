// https://leetcode.com/problems/remove-nth-node-from-end-of-list/

package main

func removeNthFromEnd(head *ListNode, n int) *ListNode {
	slow, fast := head, head

	for range n {
		fast = fast.Next
	}

	if fast == nil {
		return head.Next
	}

	for fast.Next != nil {
		fast = fast.Next
		slow = slow.Next
	}

	slow.Next = slow.Next.Next
	return head
}
