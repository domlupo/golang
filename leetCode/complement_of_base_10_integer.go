// https://leetcode.com/problems/complement-of-base-10-integer/

package main

func bitwiseComplement(n int) int {
	if n == 0 {
		return 1
	}

	var maximum int
	remainder := n
	for remainder > 0 {
		maximum = maximum*2 + 1
		remainder /= 2
	}
	return maximum - n
}
