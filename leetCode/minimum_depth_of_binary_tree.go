// https://leetcode.com/problems/minimum-depth-of-binary-tree/

package main

// iterative
func minDepth(root *TreeNode) int {
	m := 0
	if root == nil {
		return m
	}

	nodes := []*TreeNode{root}
	depths := []int{1}
	for len(nodes) != 0 {
		curr := nodes[0]
		depth := depths[0]
		m = depth

		nodes = nodes[1:]
		depths = depths[1:]

		if curr.Left == nil && curr.Right == nil {
			break
		}
		if curr.Left != nil {
			nodes = append(nodes, curr.Left)
			depths = append(depths, depth+1)
		}
		if curr.Right != nil {
			nodes = append(nodes, curr.Right)
			depths = append(depths, depth+1)
		}
	}
	return m
}

// recursive
func minDepthRecursive(root *TreeNode) int {
	if root == nil {
		return 0
	}

	l := minDepthRecursive(root.Left)
	r := minDepthRecursive(root.Right)

	if l == 0 {
		return r + 1
	} else if r == 0 {
		return l + 1
	} else {
		return min(l, r) + 1
	}
}
