// https://leetcode.com/problems/shortest-word-distance/

package main

import "math"

func shortestDistance(wordsDict []string, word1 string, word2 string) int {
	s := math.MaxInt

	p1 := -1
	p2 := -1
	for i, w := range wordsDict {
		if w == word1 {
			p1 = i
		}
		if w == word2 {
			p2 = i
		}
		if p1 == -1 || p2 == -1 {
			continue
		}
		if p1 > p2 {
			s = min(s, p1-p2)
		}
		if p2 > p1 {
			s = min(s, p2-p1)
		}
	}
	return s
}
