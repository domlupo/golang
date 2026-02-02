// https://leetcode.com/problems/binary-tree-preorder-traversal/

package main

// iterative
func preorderTraversal(root *TreeNode) []int {
	var t []int

	if root == nil {
		return t
	}

	s := []*TreeNode{root}
	for len(s) != 0 {
		curr := s[len(s)-1]
		s = s[:len(s)-1]
		t = append(t, curr.Val)
		if curr.Right != nil {
			s = append(s, curr.Right)
		}
		if curr.Left != nil {
			s = append(s, curr.Left)
		}
	}

	return t
}

// recursive
func preorderTraversalRecursive(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}
	nodes := []int{root.Val}
	nodes = append(nodes, preorderTraversalRecursive(root.Left)...)
	nodes = append(nodes, preorderTraversalRecursive(root.Right)...)
	return nodes
}
