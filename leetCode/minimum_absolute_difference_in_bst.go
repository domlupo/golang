// https://leetcode.com/problems/minimum-absolute-difference-in-bst/

package main

func getMinimumDifference(root *TreeNode) int {
	// don't use math.MaxInt due to weirdness near bounds
	// i := -math.MaxInt
	// i -= 1
	// i *= -1 is still a negative number
	minVal := 1_000_000
	maxVal := -1_000_000
	minDiff := 1_000_000
	h1(root, &minVal, &maxVal, &minDiff)
	return minDiff
}

func h1(n *TreeNode, minVal, maxVal, minDiff *int) {
	if n == nil {
		return
	}

	h1(n.Left, minVal, maxVal, minDiff)

	l := n.Val - *minVal
	u := *maxVal - n.Val

	if l < 0 {
		l *= -1
	}
	if u < 0 {
		u *= -1
	}

	*minDiff = min(*minDiff, l)
	*minDiff = min(*minDiff, u)

	if n.Val < *minVal {
		*minVal = n.Val
	}
	if n.Val > *maxVal {
		*maxVal = n.Val
	}

	h1(n.Right, minVal, maxVal, minDiff)
}
