// https://leetcode.com/problems/find-common-characters/

package main

import "math"

func commonChars(words []string) []string {
	counts := make([]int, 26)
	for i := range counts {
		counts[i] = math.MaxInt
	}

	for _, w := range words {
		wordCounts := make([]int, 26)
		for _, c := range w {
			wordCounts[c-'a']++
		}

		for i, c := range counts {
			counts[i] = min(c, wordCounts[i])
		}
	}

	var ans []string
	for i, count := range counts {
		for count > 0 {
			ans = append(ans, string('a'+i))
			count--
		}
	}
	return ans
}
