// https://leetcode.com/problems/linked-list-cycle/
// N is the number of nodes
// time complexity is O(N)
// space complexity is O(1)

package main

func hasCycle(head *ListNode) bool {
	slow := head
	fast := head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast {
			return true
		}
	}
	return false
}
