// https://leetcode.com/problems/closest-binary-search-tree-value/

package main

func closestValue(root *TreeNode, target float64) int {
	closest := root.Val

	s := []*TreeNode{root}
	for len(s) > 0 {
		curr := s[0]
		s = s[1:]

		diff1 := float64(curr.Val) - target
		diff2 := float64(closest) - target
		if diff1 < 0 {
			diff1 *= -1
		}
		if diff2 < 0 {
			diff2 *= -1
		}
		if diff1 == diff2 {
			closest = min(closest, curr.Val)
		} else if diff1 < diff2 {
			closest = curr.Val
		}

		if curr.Left != nil {
			s = append(s, curr.Left)
		}
		if curr.Right != nil {
			s = append(s, curr.Right)
		}
	}

	return closest
}
