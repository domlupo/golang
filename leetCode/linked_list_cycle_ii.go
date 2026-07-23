// https://leetcode.com/problems/linked-list-cycle-ii/
// N is the number of nodes
// time complexity is O(N)
// space complexity is O(1)

package main

func detectCycle(head *ListNode) *ListNode {
	slow, fast := head, head
	var cycle bool
	for fast != nil && fast.Next != nil {
		fast = fast.Next.Next
		slow = slow.Next
		if fast == slow {
			cycle = true
			break
		}
	}

	if cycle {
		fast = head
		for fast != slow {
			fast = fast.Next
			slow = slow.Next
		}
		return fast
	}

	return nil
}
