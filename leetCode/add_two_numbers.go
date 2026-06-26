// https://leetcode.com/problems/add-two-numbers

package main

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	dummy := &ListNode{}
	var carry int
	curr := dummy

	for l1 != nil || l2 != nil || carry == 1 {
		var sum int
		if l1 != nil {
			sum += l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			sum += l2.Val
			l2 = l2.Next
		}
		sum += carry

		carry = sum / 10
		n := &ListNode{
			Val: sum % 10,
		}
		curr.Next = n
		curr = n
	}

	return dummy.Next
}
