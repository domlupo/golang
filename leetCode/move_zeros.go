// https://leetcode.com/problems/move-zeroes/

package main

func moveZeroes(nums []int) {
	for i := len(nums); i >= 1; i-- {
		if i < len(nums) && nums[i-1] == 0 && nums[i] != 0 {
			nums[i-1], nums[i] = nums[i], nums[i-1]
			i += 2
		}
	}
}
