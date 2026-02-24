// https://leetcode.com/problems/sort-array-by-parity/

package main

func sortArrayByParity(nums []int) []int {
	var p1 int

	for p2, n := range nums {
		if n%2 == 0 {
			nums[p1], nums[p2] = nums[p2], nums[p1]
			p1++
		}
	}

	return nums
}
