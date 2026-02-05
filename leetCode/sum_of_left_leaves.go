// https://leetcode.com/problems/sum-of-left-leaves/

package main

func sumOfLeftLeaves(root *TreeNode) int {
	return sums(root, false)
}

func sums(root *TreeNode, isLeft bool) int {
	if root.Left == nil && root.Right == nil && isLeft {
		return root.Val
	}
	if root.Left == nil && root.Right == nil && !isLeft {
		return 0
	}
	if root.Left == nil {
		return sums(root.Right, false)
	}
	if root.Right == nil {
		return sums(root.Left, true)
	}
	return sums(root.Left, true) + sums(root.Right, false)
}
