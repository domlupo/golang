// https://leetcode.com/problems/first-unique-character-in-a-string/

package main

func firstUniqChar(s string) int {
	counts := make(map[rune]int)

	for _, l := range s {
		counts[l] = counts[l] + 1
	}

	for i, l := range s {
		if counts[l] == 1 {
			return i
		}
	}

	return -1
}
