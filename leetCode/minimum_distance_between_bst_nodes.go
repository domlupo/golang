// https://leetcode.com/problems/minimum-distance-between-bst-nodes/

package main

func minDiffInBST(root *TreeNode) int {
	var sorted []int

	curr := root
	var s []*TreeNode
	for curr != nil || len(s) > 0 {
		for curr != nil {
			s = append(s, curr)
			curr = curr.Left
		}
		curr = s[len(s)-1]
		s = s[:len(s)-1]
		sorted = append(sorted, curr.Val)
		curr = curr.Right
	}

	m := sorted[1] - sorted[0]
	for i := 1; i < len(sorted)-1; i++ {
		m = min(m, sorted[i+1]-sorted[i])
	}
	return m
}
