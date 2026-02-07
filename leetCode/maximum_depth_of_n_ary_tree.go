// https://leetcode.com/problems/maximum-depth-of-n-ary-tree/

package main

type Node struct {
	Val      int
	Children []*Node
}

func maxDepth1(root *Node) int {
	if root == nil {
		return 0
	}

	m := 0
	for _, c := range root.Children {
		m = max(m, maxDepth1(c))
	}

	return m + 1
}
