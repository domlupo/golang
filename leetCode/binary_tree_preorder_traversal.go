// https://leetcode.com/problems/binary-tree-preorder-traversal/

package main

// iterative
func preorderTraversal(root *TreeNode) []int {
	var t []int
	s := []*TreeNode{root}

	for len(s) > 0 {
		curr := s[len(s)-1]
		s = s[:len(s)-1]
		if curr == nil {
			continue
		}
		t = append(t, curr.Val)
		s = append(s, curr.Right, curr.Left)
	}

	return t
}

// recursive
func preorderTraversalRecursive(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}
	v := []int{root.Val}
	v = append(v, preorderTraversalRecursive(root.Left)...)
	v = append(v, preorderTraversalRecursive(root.Right)...)
	return v
}
