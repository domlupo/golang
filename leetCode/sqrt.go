// https://leetcode.com/problems/sqrtx/

package main

func mySqrt(x int) int {
	var l int
	r := x

	for l <= r {
		m := l + (r-l)/2
		if m*m > x {
			r = m - 1
		} else {
			l = m + 1
		}
	}

	return l - 1
}
