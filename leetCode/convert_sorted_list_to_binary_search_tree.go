// https://leetcode.com/problems/convert-sorted-list-to-binary-search-tree/
// N is the number of nodes in list
// time complexity is O(N*log(N)). N from finding middle node.
// log(N) from splitting list in half each time.
// space complexity is O(log(N)) from recursion stack since it is a height balanced tree.
// space complexity would be O(N) for a skewed tree.
// There is also a simple solution with O(N) time and O(N) space.
// The solution is just storing node values in an array.
// There is a more complex solution with O(N) time and O(log(N)) space.

package main

func sortedListToBST(head *ListNode) *TreeNode {
	if head == nil {
		return nil
	}
	m := splitMiddleNode(head)
	n := &TreeNode{
		Val: m.Val,
	}
	if head == m {
		return n
	}
	n.Right = sortedListToBST(m.Next)
	n.Left = sortedListToBST(head)
	return n
}

func splitMiddleNode(head *ListNode) *ListNode {
	var prev *ListNode
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		prev = slow
		slow = slow.Next
		fast = fast.Next.Next
	}
	if prev != nil {
		prev.Next = nil
	}

	return slow
}
