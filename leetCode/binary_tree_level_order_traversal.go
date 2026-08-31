// https://leetcode.com/problems/binary-tree-level-order-traversal/
// N is the number of tree nodes
// time complexity is O(N)
// space complexity is O(N)

package main

func levelOrder(root *TreeNode) [][]int {
	var ans [][]int
	travel128(root, 0, &ans)
	return ans
}

func travel128(n *TreeNode, depth int, ans *[][]int) {
	if n == nil {
		return
	}
	if depth == len(*ans) {
		*ans = append(*ans, []int{})
	}
	(*ans)[depth] = append((*ans)[depth], n.Val)
	travel128(n.Left, depth+1, ans)
	travel128(n.Right, depth+1, ans)
}
