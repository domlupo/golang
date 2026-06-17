// https://leetcode.com/problems/convert-sorted-array-to-binary-search-tree/

package main

func sortedArrayToBST(nums []int) *TreeNode {
	if len(nums) == 0 {
		return nil
	}
	m := len(nums) / 2
	n := &TreeNode{
		Val:   nums[m],
		Left:  sortedArrayToBST(nums[:m]),
		Right: sortedArrayToBST(nums[m+1:]),
	}
	return n
}
