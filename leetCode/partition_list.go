// https://leetcode.com/problems/partition-list/
// N is the length of the linked list
// time complexity is O(N)
// space complexity is O(1)

package main

func partition(head *ListNode, x int) *ListNode {
	first := &ListNode{}
	second := &ListNode{}
	firstHead, secondHead := first, second
	curr := head

	for curr != nil {
		if curr.Val < x {
			first.Next = curr
			first = first.Next
		} else {
			second.Next = curr
			second = second.Next
		}

		curr = curr.Next
	}

	if first == firstHead {
		return secondHead.Next
	} else if second == secondHead {
		return firstHead.Next
	} else {
		first.Next = secondHead.Next
		second.Next = nil
		return firstHead.Next
	}
}
