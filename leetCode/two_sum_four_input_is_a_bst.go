// https://leetcode.com/problems/two-sum-iv-input-is-a-bst/

package main

func findTarget(root *TreeNode, k int) bool {
	var sorted []int
	h4(root, &sorted)

	var p1 int
	p2 := len(sorted) - 1

	for p1 < p2 {
		sum := sorted[p1] + sorted[p2]
		if sum == k {
			return true
		} else if sum < k {
			p1++
		} else {
			p2--
		}
	}

	return false
}

func h4(n *TreeNode, sorted *[]int) {
	if n == nil {
		return
	}

	h4(n.Left, sorted)
	*sorted = append(*sorted, n.Val)
	h4(n.Right, sorted)
}
