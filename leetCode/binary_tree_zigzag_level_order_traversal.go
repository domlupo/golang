// https://leetcode.com/problems/binary-tree-zigzag-level-order-traversal/
// N is the number of nodes
// time complexity is O(N)
// space complexity is O(N)

package main

func zigzagLevelOrder(root *TreeNode) [][]int {
	var ans [][]int
	queue := []*TreeNode{root}
	depths := []int{1}

	for len(queue) != 0 {
		n := queue[0]
		d := depths[0]
		queue = queue[1:]
		depths = depths[1:]
		if n == nil {
			continue
		}

		if d > len(ans) {
			ans = append(ans, []int{})
		}
		ans[d-1] = append(ans[d-1], n.Val)
		queue = append(queue, n.Left, n.Right)
		depths = append(depths, d+1, d+1)
	}

	for i := 1; i < len(ans); i += 2 {
		var l int
		r := len(ans[i]) - 1
		for l < r {
			ans[i][l], ans[i][r] = ans[i][r], ans[i][l]
			l++
			r--
		}
	}

	return ans
}
