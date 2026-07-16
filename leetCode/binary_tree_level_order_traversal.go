// https://leetcode.com/problems/binary-tree-level-order-traversal/
// N is the number of tree nodes
// time complexity is O(N)
// space complexity is O(N)

package main

func levelOrder(root *TreeNode) [][]int {
	var ans [][]int
	depths := []int{1}
	queue := []*TreeNode{root}

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

	return ans
}
