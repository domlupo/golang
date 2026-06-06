// https://leetcode.com/problems/remove-element/description/

package main

func removeElement(nums []int, val int) int {
	var index int

	for i := range nums {
		if nums[i] != val {
			nums[index] = nums[i]
			index++
		}
	}

	return index
}
