// https://leetcode.com/problems/middle-of-the-linked-list/

package main

func middleNode(head *ListNode) *ListNode {
	p1, p2 := head, head

	for p2.Next != nil && p2.Next.Next != nil {
		p1 = p1.Next
		p2 = p2.Next.Next
	}

	if p2.Next == nil {
		return p1
	} else {
		return p1.Next
	}
}
