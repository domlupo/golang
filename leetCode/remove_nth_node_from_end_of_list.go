// https://leetcode.com/problems/remove-nth-node-from-end-of-list/
// N is the number of nodes
// time complexity is O(N)
// space complexity is O(1)

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
