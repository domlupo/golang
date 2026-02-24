// https://leetcode.com/problems/n-repeated-element-in-size-2n-array/

package main

func repeatedNTimes(nums []int) int {
	for i := 0; i < len(nums)-2; i++ {
		if nums[i] == nums[i+1] || nums[i] == nums[i+2] || nums[i] == nums[i+3] {
			return nums[i]
		} else if nums[i+1] == nums[i+2] || nums[i+1] == nums[i+3] {
			return nums[i+1]
		} else if nums[i+2] == nums[i+3] {
			return nums[i+2]
		}
	}
	panic("no solution found")
}
