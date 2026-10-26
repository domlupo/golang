// https://leetcode.com/problems/diameter-of-binary-tree/
// N is the length of nodes, H is the height of the tree.
// time complexity is O(N)
// space complexity is O(H), which is O(N) if skewed tree.

package main

func diameterOfBinaryTree(root *TreeNode) int {
	var longest int
	maxDepth543(root, &longest)
	return longest
}

func maxDepth543(n *TreeNode, diameter *int) int {
	if n == nil {
		return 0
	}
	l := maxDepth543(n.Left, diameter)
	r := maxDepth543(n.Right, diameter)
	*diameter = max(*diameter, l+r)
	return 1 + max(l, r)
}
