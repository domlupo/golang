// https://leetcode.com/problems/binary-tree-tilt/

package main

func findTilt(root *TreeNode) int {
	var tiltSum int
	treeSum(root, &tiltSum)
	return tiltSum
}

func treeSum(n *TreeNode, tiltSum *int) int {
	if n == nil {
		return 0
	}

	l := treeSum(n.Left, tiltSum)
	r := treeSum(n.Right, tiltSum)
	diff := l - r
	if diff < 0 {
		diff *= -1
	}
	*tiltSum += diff

	return n.Val + l + r
}
