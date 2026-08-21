// https://leetcode.com/problems/reverse-nodes-in-k-group/
// N is the length of nodes
// time complexity is O(N)
// space complexity is O(1)

package main

func reverseKGroup(head *ListNode, k int) *ListNode {
	dummy := &ListNode{Next: head}
	prevGroupEnd := dummy

	for {
		kth := getKthNode(prevGroupEnd, k)
		if kth == nil {
			break
		}

		prev, curr := kth.Next, prevGroupEnd.Next
		for range k {
			next := curr.Next
			curr.Next = prev
			prev, curr = curr, next
		}

		save := prevGroupEnd.Next
		prevGroupEnd.Next = kth
		prevGroupEnd = save
	}

	return dummy.Next
}

func getKthNode(curr *ListNode, k int) *ListNode {
	for curr != nil && k > 0 {
		curr = curr.Next
		k--
	}
	return curr
}
