// https://leetcode.com/problems/number-of-segments-in-a-string/

package main

func countSegments(s string) int {
	var count int

	for i := 0; i < len(s); i++ {
		if s[i] != ' ' {
			count++
			i++
		}
		for i < len(s) && s[i] != ' ' {
			i++
		}
	}

	return count
}
