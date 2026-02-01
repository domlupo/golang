// https://leetcode.com/problems/convert-sorted-array-to-binary-search-tree/

package main

func sortedArrayToBST(nums []int) *TreeNode {
	if len(nums) == 0 {
		return nil
	}
	mid := len(nums) / 2
	n := &TreeNode{
		Val:   nums[mid],
		Left:  sortedArrayToBST(nums[0:mid]),
		Right: sortedArrayToBST(nums[mid+1:]),
	}
	return n
}
