// https://leetcode.com/problems/binary-tree-paths/

package main

import "strconv"

func binaryTreePaths(root *TreeNode) []string {
	var ans []string
	paths := []string{strconv.Itoa(root.Val)}
	nodes := []*TreeNode{root}

	for len(nodes) > 0 {
		n := nodes[0]
		p := paths[0]
		nodes = nodes[1:]
		paths = paths[1:]

		if n.Left == nil && n.Right == nil {
			ans = append(ans, p)
		}
		if n.Right != nil {
			nodes = append(nodes, n.Right)
			p := p + "->" + strconv.Itoa(n.Right.Val)
			paths = append(paths, p)
		}
		if n.Left != nil {
			nodes = append(nodes, n.Left)
			p := p + "->" + strconv.Itoa(n.Left.Val)
			paths = append(paths, p)
		}
	}

	return ans
}
