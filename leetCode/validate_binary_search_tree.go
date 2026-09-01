// https://leetcode.com/problems/validate-binary-search-tree/
// N is the number of nodes in tree
// time complexity is O(N)
// space complexity is O(N)

package main

import "math"

func isValidBST(root *TreeNode) bool {
	return valid98(root, math.MinInt, math.MaxInt)
}

func valid98(n *TreeNode, minimum, maximum int) bool {
	if n == nil {
		return true
	}
	if n.Val >= maximum || n.Val <= minimum {
		return false
	}
	l := valid98(n.Left, minimum, n.Val)
	r := valid98(n.Right, n.Val, maximum)
	return l && r
}
