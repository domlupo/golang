// https://leetcode.com/problems/arranging-coins/

package main

func arrangeCoins(n int) int {
	r := 1
	for n >= r {
		n -= r
		r++
	}
	return r - 1
}
