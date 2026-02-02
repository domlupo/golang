// https://leetcode.com/problems/path-sum/

package main

// recursive (same runtime, better memory)
func hasPathSumRecursive(root *TreeNode, targetSum int) bool {
	if root == nil {
		return false
	}

	if root.Left == nil && root.Right == nil && root.Val == targetSum {
		return true
	}

	l := hasPathSumRecursive(root.Left, targetSum-root.Val)
	r := hasPathSumRecursive(root.Right, targetSum-root.Val)

	return l || r
}

// iterative
func hasPathSum(root *TreeNode, targetSum int) bool {
	if root == nil {
		return false
	}

	stack := []*TreeNode{root}
	sums := []int{0}
	for len(stack) > 0 {
		curr := stack[0]
		sum := sums[0] + curr.Val
		if curr.Left == nil && curr.Right == nil && sum == targetSum {
			return true
		}
		stack = stack[1:]
		sums = sums[1:]

		if curr.Left != nil {
			stack = append(stack, curr.Left)
			sums = append(sums, sum)
		}
		if curr.Right != nil {
			stack = append(stack, curr.Right)
			sums = append(sums, sum)
		}
	}

	return false
}
