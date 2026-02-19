// https://leetcode.com/problems/binary-gap/

package main

import "math"

func binaryGap(n int) int {
	gap := -math.MaxInt
	maxGap := -math.MaxInt

	for n > 0 {
		if n%2 == 0 {
			gap++
		} else {
			maxGap = max(maxGap, gap)
			gap = 1
		}
		n /= 2
	}

	return max(0, maxGap)
}
