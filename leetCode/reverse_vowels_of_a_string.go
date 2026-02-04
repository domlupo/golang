// https://leetcode.com/problems/reverse-vowels-of-a-string/

package main

func reverseVowels(s string) string {
	var runes []rune
	for _, s := range s {
		runes = append(runes, s)
	}

	vowels := map[rune]bool{
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

	l := 0
	r := len(runes) - 1

	for l < r {
		_, lOk := vowels[runes[l]]
		if !lOk {
			l++
			continue
		}
		_, rOk := vowels[runes[r]]
		if !rOk {
			r--
			continue
		}
		runes[l], runes[r] = runes[r], runes[l]
		l++
		r--
	}
	return string(runes)
}
