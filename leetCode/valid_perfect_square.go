// https://leetcode.com/problems/valid-perfect-square/

package main

func isPerfectSquare(num int) bool {
	l := 0
	r := num

	for l < r {
		m := l + (r-l)/2
		if m*m == num {
			return true
		} else if m*m < num {
			l = m + 1
		} else {
			r = m - 1
		}
	}
	return l*l == num
}
