// https://leetcode.com/problems/maximum-depth-of-binary-tree/

package main

func maxDepth(root *TreeNode) int {
	q := []*TreeNode{root}
	depths := []int{1}
	var depth int

	for len(q) > 0 {
		curr := q[0]
		d := depths[0]
		q = q[1:]
		depths = depths[1:]
		if curr != nil {
			depth = max(depth, d)
			q = append(q, curr.Left, curr.Right)
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
