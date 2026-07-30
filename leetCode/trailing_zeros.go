// https://leetcode.com/problems/factorial-trailing-zeroes/
// N is the passed value
// time complexity is O(log(N))
// space complexity is O(1)

package main

func trailingZeroes(n int) int {
	var zeros int
	multiple := 5
	for n >= multiple {
		zeros += (n / multiple)
		multiple *= 5
	}
	return zeros
}
