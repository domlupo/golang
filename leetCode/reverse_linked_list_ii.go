// https://leetcode.com/problems/reverse-linked-list-ii/
// N is the number of nodes
// time complexity is O(N)
// space complexity is O(1)

package main

func reverseBetween(head *ListNode, left int, right int) *ListNode {
	// this base case helps avoid a lot of issues. For example, a
	// single node, left = 1 and right = 1. It is very easy for the
	// dummy node to be linked to the final result and requires complex
	// logic to avoid.
	if left == right {
		return head
	}

	dummy := &ListNode{Next: head}
	var leftPrev, leftNode *ListNode
	count := 1
	prev := dummy
	curr := head

	for curr != nil {
		if count == left {
			leftPrev = prev
			leftNode = curr
		}

		if count == right {
			leftPrev.Next = curr
			leftNode.Next = curr.Next
		}

		if count >= left && count <= right {
			save := curr.Next
			curr.Next = prev
			prev, curr = curr, save
		} else {
			prev, curr = curr, curr.Next
		}

		count++
	}

	return dummy.Next
}
