// https://leetcode.com/problems/power-of-two/

package main

func isPowerOfTwo(n int) bool {
	// all powers of two have only 1 bit set.
	// (n&(n-1)) equals 0 when only 1 bit is set.
	return n > 0 && (n&(n-1)) == 0
}
