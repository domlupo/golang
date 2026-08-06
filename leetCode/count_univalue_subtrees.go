// https://leetcode.com/problems/count-univalue-subtrees/
// N is the number of nodes
// time complexity is O(N)
// space complexity is O(N)

package main

func countUnivalSubtrees(root *TreeNode) int {
	if root == nil {
		return 0
	}

	var count int
	helper(root, root.Val, &count)
	return count
}

func helper(n *TreeNode, val int, count *int) bool {
	if n == nil {
		return true
	}

	left := helper(n.Left, n.Val, count)
	right := helper(n.Right, n.Val, count)

	if !left || !right {
		return false
	}

	*count++

	return n.Val == val
}
