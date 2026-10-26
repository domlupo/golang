// https://leetcode.com/problems/sqrtx/
// N is the value of X
// time complexity is O(log(N))
// space complexity is O(1)

package main

func mySqrt(x int) int {
	l, r := 1, x
	for l <= r {
		m := (l + r) / 2
		if m*m > x {
			r = m - 1
		} else {
			l = m + 1
		}
	}
	return r
}
