// https://leetcode.com/problems/leaf-similar-trees/

package main

func leafSimilar(root1 *TreeNode, root2 *TreeNode) bool {
	var root1Vals, root2Vals []int
	leafVals(root1, &root1Vals)
	leafVals(root2, &root2Vals)

	if len(root1Vals) != len(root2Vals) {
		return false
	}

	for i := 0; i < len(root1Vals); i++ {
		if root1Vals[i] != root2Vals[i] {
			return false
		}
	}

	return true
}

func leafVals(node *TreeNode, vals *[]int) {
	if node == nil {
		return
	}

	leafVals(node.Left, vals)
	if node.Left == nil && node.Right == nil {
		*vals = append(*vals, node.Val)
	}
	leafVals(node.Right, vals)
}
