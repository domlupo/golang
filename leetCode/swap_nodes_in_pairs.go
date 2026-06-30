// https://leetcode.com/problems/swap-nodes-in-pairs/

package main

func swapPairs(head *ListNode) *ListNode {
	dummy := &ListNode{
		Next: head,
	}
	prev := dummy
	curr := prev.Next

	for curr != nil && curr.Next != nil {
		next := curr.Next
		curr.Next = curr.Next.Next
		next.Next = curr
		prev.Next = next

		prev = curr
		curr = prev.Next
	}

	return dummy.Next
}

func swapPairsRecursive(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	curr := head
	next := head.Next
	curr.Next = swapPairs(next.Next)
	next.Next = curr
	return next
}
