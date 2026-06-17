// https://leetcode.com/problems/palindrome-linked-list/

package main

func isPalindromeLinkedList(head *ListNode) bool {
	// set slow to middle node
	slow := head
	fast := head
	for fast != nil && fast.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next
	}

	// reverse iterative (less space than recursive)
	var prev *ListNode
	curr := slow
	for curr != nil {
		save := curr.Next
		curr.Next = prev
		prev, curr = curr, save
	}
	rev := prev

	// compare
	p1 := rev
	p2 := head
	for p1 != nil && p2 != nil {
		if p1.Val != p2.Val {
			return false
		}
		p1 = p1.Next
		p2 = p2.Next
	}
	return true
}
