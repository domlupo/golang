// https://leetcode.com/problems/increasing-order-search-tree/

package main

func increasingBST(root *TreeNode) *TreeNode {
	var vals []int
	travelInOrder(root, &vals)

	curr := &TreeNode{
		Val: vals[0],
	}
	root = curr
	for i := 1; i < len(vals); i++ {
		node := &TreeNode{
			Val: vals[i],
		}
		curr.Right = node
		curr = node
	}
	return root
}

func travelInOrder(node *TreeNode, vals *[]int) {
	if node == nil {
		return
	}
	travelInOrder(node.Left, vals)
	*vals = append(*vals, node.Val)
	travelInOrder(node.Right, vals)
}
