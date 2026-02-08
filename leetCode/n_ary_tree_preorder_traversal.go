// https://leetcode.com/problems/n-ary-tree-preorder-traversal/

package main

func preorderRecursive(root *Node) []int {
	var o []int
	if root == nil {
		return o
	}

	o = append(o, root.Val)
	for _, c := range root.Children {
		o = append(o, preorderRecursive(c)...)
	}

	return o
}

func preorder(root *Node) []int {
	var o []int
	if root == nil {
		return o
	}

	s := []*Node{root}
	for len(s) > 0 {
		curr := s[len(s)-1]
		s = s[:len(s)-1]
		o = append(o, curr.Val)
		for i := len(curr.Children) - 1; i >= 0; i-- {
			s = append(s, curr.Children[i])
		}
	}

	return o
}
