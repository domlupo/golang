// https://leetcode.com/problems/binary-tree-paths/

package main

import "strconv"

func binaryTreePaths(root *TreeNode) []string {
	var o []string

	p := []string{strconv.Itoa(root.Val)}
	n := []*TreeNode{root}
	for len(n) > 0 {
		currN := n[len(n)-1]
		currP := p[len(p)-1]

		n = n[:len(n)-1]
		p = p[:len(p)-1]

		if currN.Left == nil && currN.Right == nil {
			o = append(o, currP)
		}
		if currN.Right != nil {
			n = append(n, currN.Right)
			path := currP + "->" + strconv.Itoa(currN.Right.Val)
			p = append(p, path)
		}
		if currN.Left != nil {
			n = append(n, currN.Left)
			path := currP + "->" + strconv.Itoa(currN.Left.Val)
			p = append(p, path)
		}
	}

	return o
}
