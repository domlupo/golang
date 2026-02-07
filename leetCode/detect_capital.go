// https://leetcode.com/problems/detect-capital/

package main

func detectCapitalUse(word string) bool {
	var upperCount int
	var lowerCount int

	for _, l := range word {
		if isLower(l) {
			lowerCount++
		} else {
			upperCount++
		}
	}

	if len(word) == lowerCount || len(word) == upperCount {
		return true
	}
	if upperCount == 1 && len(word)-1 == lowerCount && !isLower(rune(word[0])) {
		return true
	}
	return false
}

func isLower(r rune) bool {
	return r >= 'a' && r <= 'z'
}
