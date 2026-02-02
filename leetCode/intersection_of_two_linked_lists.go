// https://leetcode.com/problems/intersection-of-two-linked-lists/

package main

func getIntersectionNode(headA, headB *ListNode) *ListNode {
	p1 := headA
	p2 := headB

	for {
		if p1 == nil && p2 == nil {
			return nil
		}
		if p1 == nil {
			p1 = headB
		}
		if p2 == nil {
			p2 = headA
		}
		if p1 == p2 {
			return p1
		}
		p1 = p1.Next
		p2 = p2.Next
	}
}
