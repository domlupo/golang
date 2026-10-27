// https://leetcode.com/problems/convert-binary-number-in-a-linked-list-to-integer/
// N is the number of nodes
// time complexity is O(N)
// space complexity is O(1)

package main

func getDecimalValue(head *ListNode) int {
	var ans int
	curr := head
	for curr != nil {
		ans = ans*2 + curr.Val
		curr = curr.Next
	}
	return ans
}
