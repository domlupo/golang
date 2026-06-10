// https://leetcode.com/problems/maximum-depth-of-binary-tree/

package main

// iteration skipped
func maxDepth(root *TreeNode) int {
	stack := []*TreeNode{root}
	depths := []int{1}
	var depth int

	for len(stack) > 0 {
		curr := stack[0]
		d := depths[0]
		stack = stack[1:]
		depths = depths[1:]
		if curr != nil {
			depth = max(depth, d)
			stack = append(stack, curr.Left, curr.Right)
			depths = append(depths, d+1, d+1)
		}
	}

	return depth
}

// recursion
func maxDepthRecursive(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return 1 + max(maxDepth(root.Left), maxDepth(root.Right))
}
