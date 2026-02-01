// https://leetcode.com/problems/binary-tree-inorder-traversal/

package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// recursive solution

func inorderTraversalRecursive(root *TreeNode) []int {
	var values []int
	if root == nil {
		return values
	}
	values = append(values, inorderTraversalRecursive(root.Left)...)
	values = append(values, root.Val)
	values = append(values, inorderTraversalRecursive(root.Right)...)
	return values
}

// iterative solution

func inorderTraversal(root *TreeNode) []int {
	var res []int
	var stack []*TreeNode

	curr := root
	for curr != nil || len(stack) != 0 {
		for curr != nil {
			stack = append(stack, curr)
			curr = curr.Left
		}
		curr = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		res = append(res, curr.Val)
		curr = curr.Right
	}

	return res
}
