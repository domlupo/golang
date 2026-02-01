// https://leetcode.com/problems/balanced-binary-tree/

package main

func isBalanced(root *TreeNode) bool {
	if root == nil {
		return true
	}

	l := depth(root.Left)
	r := depth(root.Right)
	diff := l - r
	if diff >= -1 && diff <= 1 {
		return isBalanced(root.Left) && isBalanced(root.Right)
	} else {
		return false
	}
}

func depth(n *TreeNode) int {
	if n == nil {
		return 0
	}
	l := depth(n.Left)
	r := depth(n.Right)
	return max(l, r) + 1
}
