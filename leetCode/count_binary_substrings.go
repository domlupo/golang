// https://leetcode.com/problems/count-binary-substrings/

package main

func countBinarySubstrings(s string) int {
	var total, prev int

	for i := 0; i < len(s); i++ {
		curr := 1
		for i+1 < len(s) && s[i] == s[i+1] {
			curr++
			i++
		}
		total += min(prev, curr)
		prev = curr
	}

	return total
}
