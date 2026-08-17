// https://leetcode.com/problems/reorder-list/
// N is the number of nodes in list
// time complexity is O(N)
// space complexity is O(1)

package main

func reorderList(head *ListNode) {
	// split list in half
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next
	}

	// reverse right half
	var prev *ListNode
	curr := slow
	for curr != nil {
		next := curr.Next
		curr.Next = prev
		prev, curr = curr, next
	}

	// merge left and reversed right
	p1, p2 := head, prev
	for p2.Next != nil {
		s1, s2 := p1.Next, p2.Next
		p1.Next, p2.Next = p2, p1.Next
		p1, p2 = s1, s2
	}
}
