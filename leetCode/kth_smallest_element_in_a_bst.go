// https://leetcode.com/problems/kth-smallest-element-in-a-bst/
// H is height of tree, K is passed k
// time complexity is O(H + K)
// space complexity is O(H)

package main

func kthSmallest(root *TreeNode, k int) int {
	n := root
	var stack []*TreeNode
	for n != nil || len(stack) > 0 {
		for n != nil {
			stack = append(stack, n)
			n = n.Left
		}
		n = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if k == 1 {
			return n.Val
		}
		k--
		n = n.Right
	}
	panic("no answer found")
}
