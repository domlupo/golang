// https://leetcode.com/problems/count-vowel-substrings-of-a-string/
// N is the length of word
// time complexity is O(N)
// space complexity is O(1)

package main

func countVowelSubstrings(word string) int {
	a, e, i, o, u := -1, -1, -1, -1, -1
	lastNonVowel := -1
	var ans int
	for ci, c := range word {
		switch c {
		case 'a':
			a = ci
		case 'e':
			e = ci
		case 'i':
			i = ci
		case 'o':
			o = ci
		case 'u':
			u = ci
		default:
			lastNonVowel = ci
		}
		ans += max(0, min(a, e, i, o, u)-lastNonVowel)
	}
	return ans
}
