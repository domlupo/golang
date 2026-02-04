// https://leetcode.com/problems/power-of-three/

package main

func isPowerOfThree(n int) bool {
	// 1162261467 is 3^19 aka largest power of three that fits in an int
	return n > 0 && 1162261467%n == 0
}
