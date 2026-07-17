// https://leetcode.com/problems/flatten-binary-tree-to-linked-list/
// N is the length of nodes
// time complexity is O(N)
// space complexity is O(1) due to Morris traversal

package main

func flatten(root *TreeNode) {
	if root == nil {
		return
	}

	n := root
	for n != nil {
		if n.Left != nil {
			rightmost := n.Left
			for rightmost.Right != nil {
				rightmost = rightmost.Right
			}
			rightmost.Right = n.Right
			n.Right = n.Left
			n.Left = nil
		}
		n = n.Right
	}
}

/*

O(N) space complexity solution

func flatten(root *TreeNode) {
	if root == nil {
		return
	}

	dummy := &TreeNode{}
	curr := dummy
	q := []*TreeNode{root}

	for len(q) != 0 {
		n := q[len(q)-1]
		q = q[:len(q)-1]
		if n.Right != nil {
			q = append(q, n.Right)
		}
		if n.Left != nil {
			q = append(q, n.Left)
		}
		curr.Right = n
		curr.Left = nil
		curr = n
	}
}

*/
