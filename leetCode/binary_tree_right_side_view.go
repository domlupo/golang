// https://leetcode.com/problems/binary-tree-right-side-view
// N is the number of nodes
// time complexity is O(N)
// space complexity is O(N)

package main

func rightSideView(root *TreeNode) []int {
	var ans []int
	var nextDepth int
	depths := []int{0}
	nodes := []*TreeNode{root}

	for len(nodes) != 0 {
		n := nodes[0]
		d := depths[0]
		nodes = nodes[1:]
		depths = depths[1:]
		if n == nil {
			continue
		} else if d == nextDepth {
			ans = append(ans, n.Val)
			nextDepth++
		}
		nodes = append(nodes, n.Right, n.Left)
		depths = append(depths, d+1, d+1)
	}

	return ans
}
