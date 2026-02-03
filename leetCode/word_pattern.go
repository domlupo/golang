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

	lToS := make(map[byte]string)
	sToL := make(map[string]byte)
	for i := 0; i < len(words); i++ {
		s, sOk := lToS[pattern[i]]
		l, lOk := sToL[words[i]]

		if sOk && s != words[i] {
			return false
		}
		if lOk && l != pattern[i] {
			return false
		}
		lToS[pattern[i]] = words[i]
		sToL[words[i]] = pattern[i]
	}

	return true
}
