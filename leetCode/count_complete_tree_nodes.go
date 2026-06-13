// https://leetcode.com/problems/count-complete-tree-nodes/

package main

func countNodes(root *TreeNode) int {
	if root == nil {
		return 0
	}

	l := depth1(root.Left)
	r := depth1(root.Right)

	ans := 1 << r
	if l == r {
		ans += countNodes(root.Right)
	} else {
		ans += countNodes(root.Left)
	}
	return ans
}

func depth1(n *TreeNode) int {
	if n == nil {
		return 0
	}
	return 1 + depth1(n.Left)
}

// 4 -> 2 + 1 + 1
// 5 -> 2 + 2 + 1
// 6 -> 4 + 1 + 1
// 7 -> 4 + 2 + 1
