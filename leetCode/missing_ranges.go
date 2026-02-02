// https://leetcode.com/problems/missing-ranges/

package main

func findMissingRanges(nums []int, lower int, upper int) [][]int {
	var o [][]int
	var i int

	for lower <= upper {
		if i >= len(nums) {
			o = append(o, []int{lower, upper})
			break
		}
		if nums[i] == lower {
			lower++
			i++
			continue
		}
		var right = upper
		if i < len(nums) {
			right = nums[i] - 1
		}
		o = append(o, []int{lower, right})
		lower = right + 1
	}

	return o
}
