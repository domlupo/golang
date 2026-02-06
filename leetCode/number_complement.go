// https://leetcode.com/problems/number-complement/

package main

func findComplement(num int) int {
	n := num
	var mask int
	for n > 0 {
		mask = mask*2 + 1
		n /= 2
	}
	return num ^ mask
}
