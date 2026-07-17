// https://leetcode.com/problems/binary-tree-level-order-traversal-ii/
// N is the number of nodes in tree
// time complexity is O(N)
// space complexity is O(N)

package main

func levelOrderBottom(root *TreeNode) [][]int {
	var ans [][]int
	queue := []*TreeNode{root}
	depths := []int{0}

	for len(queue) != 0 {
		n := queue[0]
		d := depths[0]
		queue = queue[1:]
		depths = depths[1:]
		if n == nil {
			continue
		}
		if len(ans) == d {
			ans = append(ans, []int{})
		}
		ans[d] = append(ans[d], n.Val)
		queue = append(queue, n.Left, n.Right)
		depths = append(depths, d+1, d+1)
	}

	var l int
	r := len(ans) - 1
	for l < r {
		ans[l], ans[r] = ans[r], ans[l]
		l++
		r--
	}

	return ans
}
