// https://leetcode.com/problems/count-complete-tree-nodes/
// N is the number of nodes, D is depth
// time complexity is O(D*D) aka O(log(N) * log(N))
// space complexity is O(1)

package main

func countNodes(root *TreeNode) int {
	var ans int

	n := root
	for n != nil {
		l := rightDepth(n.Left)
		r := rightDepth(n.Right)
		ans += 1 << l
		if l == r {
			n = n.Left
		} else {
			n = n.Right
		}
	}

	return ans
}

func rightDepth(n *TreeNode) int {
	var d int
	for n != nil {
		n = n.Right
		d++
	}
	return d
}
