// https://leetcode.com/problems/shortest-distance-to-a-character/

package main

func shortestToChar(s string, c byte) []int {
	l := 1_000_000
	r := 0

	var o []int
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			l = i
		}

		for r < len(s) && (r == i || s[r] != c) {
			r++
		}

		if r == len(s) {
			r = 1_000_000
		}

		o = append(o, min(absDiff1(l, i), absDiff1(r, i)))
	}
	return o
}

func absDiff1(v1, v2 int) int {
	if v1 > v2 {
		return v1 - v2
	} else {
		return v2 - v1
	}
}
