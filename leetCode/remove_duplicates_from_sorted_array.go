// https://leetcode.com/problems/remove-duplicates-from-sorted-array

package main

func removeDuplicates(nums []int) int {
	index := 1

	for i := 0; i < len(nums)-1; i++ {
		if nums[i] != nums[i+1] {
			nums[index] = nums[i+1]
			index++
		}
	}

	return index
}
