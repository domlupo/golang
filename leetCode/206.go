// https://leetcode.com/problems/reverse-linked-list/
// N is the number of nodes
// time complexity is O(N)
// space complexity is O(1)

package main

func reverseList(head *ListNode) *ListNode {
	var prev *ListNode
	curr := head
	for curr != nil {
		save := prev
		prev, curr = curr, curr.Next
		prev.Next = save
	}
	return prev
}

func reverseListRecursive(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	newHead := reverseListRecursive(head.Next)

	head.Next.Next = head
	head.Next = nil

	return newHead
}
