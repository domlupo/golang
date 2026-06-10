// https://leetcode.com/problems/symmetric-tree/

package main

// iteration
func isSymmetric(root *TreeNode) bool {
	stack := []*TreeNode{root.Left, root.Right}

	for len(stack) > 0 {
		l := stack[0]
		r := stack[1]
		stack = stack[2:]
		if l == nil && r == nil {
			continue
		} else if l == nil || r == nil || l.Val != r.Val {
			return false
		}
		stack = append(stack, l.Left, r.Right, l.Right, r.Left)
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
