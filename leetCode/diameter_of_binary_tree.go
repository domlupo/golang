// https://leetcode.com/problems/diameter-of-binary-tree/

package main

func diameterOfBinaryTree(root *TreeNode) int {
	var longest int
	h2(root, &longest)
	return longest
}

func h2(n *TreeNode, longest *int) int {
	if n == nil {
		return 0
	}

	var l int
	if n.Left != nil {
		l++
		l += h2(n.Left, longest)
	}
	var r int
	if n.Right != nil {
		r++
		r += h2(n.Right, longest)
	}

	*longest = max(*longest, l+r)

	return max(l, r)
}
