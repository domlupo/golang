// https://leetcode.com/problems/path-sum-ii/
// N is the number of nodes in the tree
// time complexity is O(N * N)
// N from traversing tree, N from copying combinations into final answer.
// space complexity is O(N) from recursion stack

package main

func pathSum(root *TreeNode, targetSum int) [][]int {
	var ans [][]int
	backtrack113(&ans, []int{}, targetSum, root)
	return ans
}

func backtrack113(ans *[][]int, comb []int, sum int, n *TreeNode) {
	if n == nil {
		return
	} else if n.Left == nil && n.Right == nil && sum == n.Val {
		comb = append(comb, n.Val)
		deepCopy := make([]int, len(comb))
		copy(deepCopy, comb)
		*ans = append(*ans, deepCopy)
		return
	}

	comb = append(comb, n.Val)
	backtrack113(ans, comb, sum-n.Val, n.Left)
	backtrack113(ans, comb, sum-n.Val, n.Right)
}
