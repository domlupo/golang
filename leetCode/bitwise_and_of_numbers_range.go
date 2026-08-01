// https://leetcode.com/problems/bitwise-and-of-numbers-range/
// time complexity is O(1)
// space complexity is O(1)

package main

func rangeBitwiseAnd(left int, right int) int {
	var shifts int
	for left != right {
		left >>= 1
		right >>= 1
		shifts++
	}
	return left << shifts
}
