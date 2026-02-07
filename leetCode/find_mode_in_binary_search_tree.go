// https://leetcode.com/problems/find-mode-in-binary-search-tree/

package main

import "math"

func findMode(root *TreeNode) []int {
	val := -math.MaxInt
	count := -math.MaxInt
	m := -math.MaxInt
	var o []int
	h(root, &val, &count, &m, &o)
	return o
}

func h(n *TreeNode, val, count, m *int, o *[]int) {
	if n == nil {
		return
	}

	h(n.Left, val, count, m, o)

	if n.Val == *val {
		*count++
	} else {
		*count = 1
		*val = n.Val
	}

	if *count > *m {
		*m = *count
		*o = []int{*val}
	} else if *count == *m {
		*o = append(*o, *val)
	}

	h(n.Right, val, count, m, o)
}
