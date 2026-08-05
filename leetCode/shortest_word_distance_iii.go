// https://leetcode.com/problems/shortest-word-distance-iii/
// N is the length of wordsDict
// time complexity is O(N)
// space complexity is O(1)

package main

import "math"

func shortestWordDistance(wordsDict []string, word1 string, word2 string) int {
	ans := math.MaxInt
	p1, p2 := -10000000, -10000000

	for i, w := range wordsDict {
		if word1 == word2 && w == word1 {
			p1 = p2
			p2 = i
		} else if w == word1 {
			p1 = i
		} else if w == word2 {
			p2 = i
		}

		if p2 > p1 {
			ans = min(ans, p2-p1)
		} else if p1 > p2 {
			ans = min(ans, p1-p2)
		}
	}

	return ans
}
