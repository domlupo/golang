// https://leetcode.com/problems/palindrome-linked-list/
// N is the number of nodes in linked list
// time complexity is O(N)
// space complexity is O(1)

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
		save := prev
		prev, curr = curr, curr.Next
		prev.Next = save
	}

	// compare
	p1 := prev
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
