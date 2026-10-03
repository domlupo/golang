// https://leetcode.com/problems/plus-one/
// N is the length of digits
// time complexity is O(N)
// space complexity is O(1)

package main

func plusOne(digits []int) []int {
	for i := len(digits) - 1; i >= 0; i-- {
		if digits[i] != 9 {
			digits[i]++
			return digits
		}
		digits[i] = 0
	}

	return append([]int{1}, digits...)
}
