// https://leetcode.com/problems/maximum-depth-of-binary-tree/

package main

func maxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}

	q := []*TreeNode{root}
	d := []int{1}
	var maxD int
	for len(q) > 0 {
		n := q[0]
		depth := d[0]
		q = q[1:]
		d = d[1:]
		if n != nil {
			maxD = max(maxD, depth)
			q = append(q, n.Left, n.Right)
			d = append(d, depth+1, depth+1)
		}
	}

	return maxD
}

// recursion
func maxDepthRecursive(root *TreeNode) int {
	if root == nil {
		return 0
	}
	return 1 + max(maxDepth(root.Left), maxDepth(root.Right))
}
