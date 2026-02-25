// https://leetcode.com/problems/cousins-in-binary-tree/

package main

type nodeInfo struct {
	*TreeNode
	parent *TreeNode
	depth  int
}

func isCousins(root *TreeNode, x int, y int) bool {
	var parents []*TreeNode
	var depths []int

	s := []nodeInfo{nodeInfo{TreeNode: root}}
	for len(s) > 0 && len(parents) != 2 {
		curr := s[0]
		s = s[1:]
		if curr.Left != nil {
			n := nodeInfo{
				TreeNode: curr.Left,
				parent:   curr.TreeNode,
				depth:    curr.depth + 1,
			}
			s = append(s, n)
		}
		if curr.Right != nil {
			n := nodeInfo{
				TreeNode: curr.Right,
				parent:   curr.TreeNode,
				depth:    curr.depth + 1,
			}
			s = append(s, n)
		}

		if curr.Val == x || curr.Val == y {
			parents = append(parents, curr.parent)
			depths = append(depths, curr.depth)
		}
	}

	return parents[0] != parents[1] && depths[0] == depths[1]
}
