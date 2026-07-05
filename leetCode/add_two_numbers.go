// https://leetcode.com/problems/add-two-numbers

package main

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	dummy := &ListNode{}
	curr := dummy
	var carry int

	for l1 != nil || l2 != nil || carry != 0 {
		sum := carry
		if l1 != nil {
			sum += l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			sum += l2.Val
			l2 = l2.Next
		}

		carry = sum / 10
		next := &ListNode{
			Val: sum % 10,
		}
		curr.Next = next
		curr = next
	}

	return dummy.Next
}
