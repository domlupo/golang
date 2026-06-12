// https://leetcode.com/problems/number-of-1-bits

package main

func hammingWeight(n int) int {
	var ans int
	for n > 0 {
		ans += n % 2
		n /= 2
	}
	return ans
}
