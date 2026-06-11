// https://leetcode.com/problems/binary-tree-postorder-traversal/

package main

// iterative
func postorderTraversal(root *TreeNode) []int {
	var ans []int
	s := []*TreeNode{root}

	for len(s) > 0 {
		curr := s[len(s)-1]
		s = s[:len(s)-1]
		if curr == nil {
			continue
		}
		s = append([]int{curr.Val}, ans...)
		s = append(s, curr.Left, curr.Right
	}

	return ans
}

// recursive
func postorderTraversalRecursive(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}

	var v []int
	v = append(v, postorderTraversalRecursive(root.Left)...)
	v = append(v, postorderTraversalRecursive(root.Right)...)
	v = append(v, root.Val)
	return v
}
