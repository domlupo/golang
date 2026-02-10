// https://leetcode.com/problems/binary-number-with-alternating-bits/

package main

func hasAlternatingBits(n int) bool {
	lastBit := -1
	for n > 0 {
		if n%2 == 1 && lastBit == 1 {
			return false
		} else if n%2 == 0 && lastBit == 0 {
			return false
		}

		lastBit = n % 2
		n /= 2
	}

	return true
}
