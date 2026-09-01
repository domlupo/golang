// https://leetcode.com/problems/kth-smallest-element-in-a-bst/
// H is height of tree, K is passed k
// time complexity is O(H + K)
// space complexity is O(H)

package main

func kthSmallest(root *TreeNode, k int) int {
	var stack []*TreeNode
	curr := root
	for {
		for curr != nil {
			stack = append(stack, curr)
			curr = curr.Left
		}
		curr = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if k == 1 {
			return curr.Val
		}
		k--
		curr = curr.Right
	}
}
