// https://leetcode.com/problems/verifying-an-alien-dictionary/

package main

func isAlienSorted(words []string, order string) bool {
	orderLookUp := make(map[byte]int)
	for i := 0; i < len(order); i++ {
		orderLookUp[order[i]] = i
	}

	for i := 1; i < len(words); i++ {
		if !isBefore(words[i-1], words[i], orderLookUp) {
			return false
		}
	}

	return true
}

func isBefore(s1, s2 string, order map[byte]int) bool {
	for i := 0; i < len(s1); i++ {
		if i == len(s2) {
			return false
		} else if order[s1[i]] > order[s2[i]] {
			return false
		} else if order[s1[i]] < order[s2[i]] {
			return true
		}
	}
	return true
}
