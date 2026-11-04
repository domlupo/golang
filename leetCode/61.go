// https://leetcode.com/problems/rotate-list/
// N is the number of nodes
// time complexity is O(N)
// space complexity is O(1)

package main

func rotateRight(head *ListNode, k int) *ListNode {
	if head == nil || head.Next == nil || k == 0 {
		return head
	}

	size := 1
	tail := head
	for tail.Next != nil {
		tail = tail.Next
		size++
	}
	tail.Next = head

	k = k % size
	i := size - k - 1
	curr := head
	for i > 0 {
		curr = curr.Next
		i--
	}

	newHead := curr.Next
	curr.Next = nil
	return newHead
}
