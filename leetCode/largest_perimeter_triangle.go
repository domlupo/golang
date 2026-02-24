// https://leetcode.com/problems/largest-perimeter-triangle/

package main

func largestPerimeter(nums []int) int {
	for i := 0; i < len(nums)-1; i++ {
		if i >= 0 && nums[i] < nums[i+1] {
			nums[i], nums[i+1] = nums[i+1], nums[i]
			i -= 2
		}
	}

	for i := 0; i < len(nums)-2; i++ {
		if nums[i] < nums[i+1]+nums[i+2] {
			return nums[i] + nums[i+1] + nums[i+2]
		}
	}
	return 0
}
