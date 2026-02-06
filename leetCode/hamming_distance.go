// https://leetcode.com/problems/hamming-distance/

package main

func hammingDistance(x int, y int) int {
	xor := x ^ y
	var d int
	for xor > 0 {
		if xor%2 == 1 {
			d++
		}
		xor /= 2
	}

	return d
}
