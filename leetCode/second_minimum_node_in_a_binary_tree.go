// https://leetcode.com/problems/second-minimum-node-in-a-binary-tree/

package main

import "math"

func findSecondMinimumValue(root *TreeNode) int {
	m := min(h5(root.Left, root.Val), h5(root.Right, root.Val))
	if m == math.MaxInt {
		return -1
	} else {
		return m
	}
}

func h5(n *TreeNode, filter int) int {
	if n == nil {
		return math.MaxInt
	}

	m := math.MaxInt
	m = min(m, h5(n.Left, filter))
	m = min(m, h5(n.Right, filter))

	if n.Val != filter {
		m = min(m, n.Val)
	}

	return m
}
