// https://leetcode.com/problems/happy-number/

package main

func isHappy(n int) bool {
	seen := make(map[int]bool)
	for n != 1 {
		var sum int
		for n > 0 {
			r := n % 10
			sum += r * r
			n /= 10
		}
		n = sum
		if seen[n] {
			return false
		}
		seen[n] = true
	}
	return true
}
