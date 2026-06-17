// https://leetcode.com/problems/sqrtx/

package main

func mySqrt(x int) int {
	var l int
	r := x

	for l <= r {
		m := (l + r) / 2
		if m*m <= x {
			l = m + 1
		} else {
			r = m - 1
		}
	}

	return r
}
