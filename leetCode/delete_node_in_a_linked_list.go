// https://leetcode.com/problems/delete-node-in-a-linked-list/
// time complexity is O(1)
// space complexity is O(1)

package main

func deleteNode(node *ListNode) {
	node.Val = node.Next.Val
	node.Next = node.Next.Next
}
