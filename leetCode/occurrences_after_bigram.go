// https://leetcode.com/problems/occurrences-after-bigram/

package main

func findOcurrences(text string, first string, second string) []string {
	var ans []string
	var f, s, t string
	var next []byte

	for i := 0; i <= len(text); i++ {
		if i == len(text) || text[i] == ' ' {
			f = s
			s = t
			t = string(next)
			next = []byte{}
			if f == first && s == second {
				ans = append(ans, t)
			}
		} else {
			next = append(next, text[i])
		}
	}

	return ans
}
