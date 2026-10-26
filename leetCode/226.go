// https://leetcode.com/problems/invert-binary-tree/
// N is the number of nodes, H is the height of the tree
// time complexity is O(N)
// space complexity is O(H)

package main

func invertTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	root.Left, root.Right = invertTree(root.Right), invertTree(root.Left)
	return root
}
