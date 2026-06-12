// https://leetcode.com/problems/missing-ranges/

package main

func findMissingRanges(nums []int, lower int, upper int) [][]int {
	numbers := []int{lower - 1}
	numbers = append(numbers, nums...)
	numbers = append(numbers, upper+1)
	var ans [][]int

	for i := 1; i < len(numbers); i++ {
		l := numbers[i-1]
		r := numbers[i]
		if l+1 == r {
			continue
		}
		ans = append(ans, []int{l + 1, r - 1})
	}

	return ans
}
