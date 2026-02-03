// https://leetcode.com/problems/palindrome-permutation/

package main

func canPermutePalindrome(s string) bool {
	counts := make(map[rune]int)
	var odds int

	for _, s := range s {
		counts[s] = counts[s] + 1
		if counts[s]%2 == 1 {
			odds++
		} else {
			odds--
		}
	}

	return odds < 2
}
