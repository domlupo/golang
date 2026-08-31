// https://leetcode.com/problems/binary-tree-right-side-view
// N is the number of nodes
// time complexity is O(N)
// space complexity is O(N)

package main

func rightSideView(root *TreeNode) []int {
	var ans []int
	addRight(root, 0, &ans)
	return ans
}

func addRight(n *TreeNode, depth int, ans *[]int) {
	if n == nil {
		return
	}
	if len(*ans) == depth {
		*ans = append(*ans, n.Val)
	}
	addRight(n.Right, depth+1, ans)
	addRight(n.Left, depth+1, ans)
}
