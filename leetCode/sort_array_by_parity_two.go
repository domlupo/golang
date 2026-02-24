// https://leetcode.com/problems/sort-array-by-parity-ii/

package main

func sortArrayByParityII(nums []int) []int {
	var e int
	o := 1

	for e < len(nums) && o < len(nums) {
		if nums[e]%2 == 0 {
			e += 2
		} else if nums[o]%2 == 1 {
			o += 2
		} else {
			nums[e], nums[o] = nums[o], nums[e]
		}
	}

	return nums
}
