// https://leetcode.com/problems/lowest-common-ancestor-of-a-binary-search-tree/
// H is the height of the tree
// time complexity is O(H)
// space complexity is O(1)

package main

func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
	left := min(p.Val, q.Val)
	right := max(p.Val, q.Val)

	n := root
	for {
		if n.Val > left && n.Val < right || n.Val == left || n.Val == right {
			return n
		}
		if n.Val > right {
			n = n.Left
		} else {
			n = n.Right
		}
	}
}
