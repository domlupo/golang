// https://leetcode.com/problems/same-tree

package main

func isSameTree(p *TreeNode, q *TreeNode) bool {
	queue := []*TreeNode{p, q}

	for len(queue) > 0 {
		n1 := queue[0]
		n2 := queue[1]
		queue = queue[2:]
		if n1 == nil && n2 == nil {
			continue
		} else if n1 == nil || n2 == nil || n1.Val != n2.Val {
			return false
		}
		queue = append(queue, n1.Left, n2.Left, n1.Right, n2.Right)
	}

	return true
}

func isSameTreeRecursive(p *TreeNode, q *TreeNode) bool {
	if p == nil && q == nil {
		return true
	} else if p == nil || q == nil || p.Val != q.Val {
		return false
	}
	return isSameTreeRecursive(p.Left, q.Left) && isSameTreeRecursive(p.Right, q.Right)
}
