// https://leetcode.com/problems/buddy-strings/

package main

import "math"

func buddyStrings(s string, goal string) bool {
	if len(s) != len(goal) {
		return false
	}

	counts := make(map[byte]int)
	maxCharCount := math.MinInt

	first, second := -1, -1
	for i := 0; i < len(s); i++ {
		counts[s[i]]++
		maxCharCount = max(maxCharCount, counts[s[i]])

		if s[i] != goal[i] && first == -1 {
			first = i
		} else if s[i] != goal[i] && second == -1 {
			second = i
		} else if s[i] != goal[i] {
			return false
		}
	}

	if first == -1 && second == -1 && maxCharCount >= 2 {
		return true
	}

	if first == -1 || second == -1 {
		return false
	}

	return s[first] == goal[second] && s[second] == goal[first]
}
