// https://leetcode.com/problems/count-good-nodes-in-binary-tree/
// N is the number of nodes
// time complexity is O(N)
// space complexity is O(N)

package main

import "math"

func goodNodes(root *TreeNode) int {
	var ans int
	travel1448(root, &ans, math.MinInt)
	return ans
}

func travel1448(n *TreeNode, count *int, maxVal int) {
	if n == nil {
		return
	}
	if n.Val >= maxVal {
		*count++
	}
	maxVal = max(n.Val, maxVal)
	travel1448(n.Left, count, maxVal)
	travel1448(n.Right, count, maxVal)
}
