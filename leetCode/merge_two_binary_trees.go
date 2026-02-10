// https://leetcode.com/problems/merge-two-binary-trees/

package main

func mergeTrees(root1 *TreeNode, root2 *TreeNode) *TreeNode {
	m := &TreeNode{}
	if root1 == nil && root2 == nil {
		return nil
	}
	if root1 == nil {
		m.Val = root2.Val
		m.Left = mergeTrees(nil, root2.Left)
		m.Right = mergeTrees(nil, root2.Right)
		return m
	}
	if root2 == nil {
		m.Val = root1.Val
		m.Left = mergeTrees(nil, root1.Left)
		m.Right = mergeTrees(nil, root1.Right)
		return m
	}
	m.Val = root1.Val + root2.Val
	m.Left = mergeTrees(root1.Left, root2.Left)
	m.Right = mergeTrees(root1.Right, root2.Right)
	return m
}
