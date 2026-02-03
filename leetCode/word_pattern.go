// https://leetcode.com/problems/word-pattern/

package main

func wordPattern(pattern string, s string) bool {
	var words []string
	var runes []rune
	for _, r := range s {
		if r != ' ' {
			runes = append(runes, r)
		} else {
			words = append(words, string(runes))
			runes = []rune{}
		}
	}
	if len(runes) > 0 {
		words = append(words, string(runes))
	}

	if len(pattern) != len(words) {
		return false
	}

	letterToString := make(map[byte]string)
	stringToLetter := make(map[string]byte)

	for i := 0; i < len(words); i++ {
		s, sOk := letterToString[pattern[i]]
		l, lOk := stringToLetter[words[i]]
		if sOk && s != words[i] {
			return false
		}
		if lOk && l != pattern[i] {
			return false
		}
		letterToString[pattern[i]] = words[i]
		stringToLetter[words[i]] = pattern[i]
	}

	return true
}
