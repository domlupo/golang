// https://leetcode.com/problems/maximum-depth-of-binary-tree/

package main

// iteration skipped

// recursion
func maxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	r := maxDepth(root.Right)
	l := maxDepth(root.Left)
	return 1 + max(r, l)
}
