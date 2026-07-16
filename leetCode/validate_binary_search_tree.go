// https://leetcode.com/problems/validate-binary-search-tree/
// M is the number of nodes in tree
// time complexity is O(M)
// space complexity is O(M)

package main

import "math"

func isValidBST(root *TreeNode) bool {
	l := valid98(root.Left, math.MinInt, root.Val)
	r := valid98(root.Right, root.Val, math.MaxInt)
	return l && r
}

func valid98(n *TreeNode, minimum, maximum int) bool {
	if n == nil {
		return true
	}

	if n.Val >= maximum || n.Val <= minimum {
		return false
	}

	l := valid98(n.Left, minimum, min(n.Val, maximum))
	r := valid98(n.Right, max(minimum, n.Val), maximum)

	return l && r
}
