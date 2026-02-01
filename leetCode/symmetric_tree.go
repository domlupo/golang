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
		}
		if l == nil || r == nil {
			return false
		}
		if l.Val != r.Val {
			return false
		}
		stack = append(stack, l.Left, r.Right, l.Right, r.Left)
	}

	return true
}

// recursive
func symmetric(n1, n2 *TreeNode) bool {
	if n1 == nil && n2 == nil {
		return true
	}
	if n1 == nil || n2 == nil {
		return false
	}
	if n1.Val != n2.Val {
		return false
	}
	return symmetric(n1.Left, n2.Right) && symmetric(n1.Right, n2.Left)
}
