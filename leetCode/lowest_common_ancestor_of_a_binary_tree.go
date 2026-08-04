// https://leetcode.com/problems/lowest-common-ancestor-of-a-binary-tree/
// N is the number of nodes
// time complexity is O(N)
// space complexity is O(N)

package main

func lowestCommonAncestor236(root, p, q *TreeNode) *TreeNode {
	return find(root, p.Val, q.Val)
}

func find(n *TreeNode, v1, v2 int) *TreeNode {
	if n == nil || n.Val == v1 || n.Val == v2 {
		return n
	}

	l := find(n.Left, v1, v2)
	r := find(n.Right, v1, v2)

	if l != nil && r != nil {
		return n
	} else if (n.Val == v1 || n.Val == v2) && (l != nil || r != nil) {
		return n
	} else if l != nil {
		return l
	} else {
		return r
	}
}
