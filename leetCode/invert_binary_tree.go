// https://leetcode.com/problems/invert-binary-tree/

package main

func invertTree(root *TreeNode) *TreeNode {
	if root == nil {
		return root
	}

	root.Right, root.Left = root.Left, root.Right

	invertTree(root.Left)
	invertTree(root.Right)

	return root
}
