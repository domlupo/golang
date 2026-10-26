// https://leetcode.com/problems/reverse-vowels-of-a-string/
// N is the length of s
// time complexity is O(N)
// space complexity is O(1)

package main

func reverseVowels(s string) string {
	vowels := map[byte]bool{
		'a': true,
		'e': true,
		'i': true,
		'o': true,
		'u': true,
		'A': true,
		'E': true,
		'I': true,
		'O': true,
		'U': true,
	}
	bytes := []byte(s)
	l, r := 0, len(s)-1
	for l < r {
		lIsVowel := vowels[bytes[l]]
		rIsVowel := vowels[bytes[r]]
		if !lIsVowel {
			l++
		} else if !rIsVowel {
			r--
		} else {
			bytes[l], bytes[r] = bytes[r], bytes[l]
			l++
			r--
		}
	}
	return string(bytes)
}
