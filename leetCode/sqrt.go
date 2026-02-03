// https://leetcode.com/problems/sqrtx/

package main

func mySqrt(x int) int {
	l := 0
	r := x

	for l <= r {
		mid := l + (r-l)/2
		if mid*mid > x {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}

	return l - 1
}
