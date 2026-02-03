// https://leetcode.com/problems/move-zeroes/

package main

func moveZeroes(nums []int) {
	var l int
	for r := 0; r < len(nums); r++ {
		if nums[r] != 0 {
			nums[l], nums[r] = nums[r], nums[l]
			l++
		}
	}
}
