// https://leetcode.com/problems/symmetric-tree/

package main

// iteration
func isSymmetric(root *TreeNode) bool {
	queue := []*TreeNode{root.Left, root.Right}

	for len(queue) > 0 {
		l := queue[0]
		r := queue[1]
		queue = queue[2:]
		if l == nil && r == nil {
			continue
		} else if l == nil || r == nil || l.Val != r.Val {
			return false
		}
		queue = append(queue, l.Left, r.Right, l.Right, r.Left)
	}

	return true
}

// recursive
func sym(l, r *TreeNode) bool {
	if l == nil && r == nil {
		return true
	} else if l == nil || r == nil || l.Val != r.Val {
		return false
	}
	return sym(l.Left, r.Right) && sym(l.Right, r.Left)
}
