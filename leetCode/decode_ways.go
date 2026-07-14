// https://leetcode.com/problems/decode-ways/
// N is the length of s
// time complexity is O(N)
// space complexity is O(1)

package main

import "strconv"

func numDecodings(s string) int {
	if s[0] == '0' {
		return 0
	}

	twoBack, oneBack := 1, 1
	for i := 1; i < len(s); i++ {
		var curr int
		if s[i] != '0' {
			curr += oneBack
		}
		n, _ := strconv.Atoi(s[i-1 : i+1])
		if n >= 10 && n <= 26 {
			curr += twoBack
		}
		twoBack, oneBack = oneBack, curr
	}

	return oneBack
}
