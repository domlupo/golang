// https://leetcode.com/problems/number-of-1-bits

package main

func hammingWeight(n int) int {
	var a int
	for {
		if n%2 == 1 {
			a++
		}
		if n == 0 || n == 1 {
			break
		}
		n /= 2
	}
	return a
}
