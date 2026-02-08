// https://leetcode.com/problems/n-ary-tree-postorder-traversal/

package main

func postorder(root *Node) []int {
	var r []int
	if root == nil {
		return r
	}

	s := []*Node{root}
	for len(s) > 0 {
		curr := s[len(s)-1]
		s = s[:len(s)-1]

		r = append(r, curr.Val)
		for _, c := range curr.Children {
			s = append(s, c)
		}
	}

	var o []int
	for i := len(r) - 1; i >= 0; i-- {
		o = append(o, r[i])
	}
	return o
}

func postorderRecursive(root *Node) []int {
	var o []int
	if root == nil {
		return o
	}

	for _, c := range root.Children {
		o = append(o, postorderRecursive(c)...)
	}
	o = append(o, root.Val)

	return o
}
