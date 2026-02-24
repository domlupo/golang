// https://leetcode.com/problems/range-sum-of-bst/

package main

func rangeSumBST(root *TreeNode, low int, high int) int {
	var sum int
	travel1(root, low, high, &sum)
	return sum
}

func travel1(n *TreeNode, low, high int, sum *int) {
	if n == nil {
		return
	}

	if n.Val > low {
		travel1(n.Left, low, high, sum)
	}
	if n.Val >= low && n.Val <= high {
		*sum += n.Val
	}
	if n.Val < high {
		travel1(n.Right, low, high, sum)
	}
}
