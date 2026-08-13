// https://leetcode.com/problems/minimum-window-substring/
// S is the length of s, T is the length of t
// time complexity is O(S + T)
// space complexity is O(S + T)

package main

func minWindow(s string, t string) string {
	var ans string
	if len(t) > len(s) {
		return ans
	}

	var sFreq, tFreq [52]int
	var required int
	for _, c := range t {
		tc := convert76(byte(c))
		tFreq[tc]++
		if tFreq[tc] == 1 {
			required++
		}
	}

	var created, l int
	minWindow := len(s) + 1
	for r := 0; r < len(s); r++ {
		rc := convert76(s[r])
		sFreq[rc]++
		if sFreq[rc] == tFreq[rc] {
			created++
		}

		for created == required {
			str := s[l : r+1]
			if len(str) < minWindow {
				ans = str
				minWindow = len(ans)
			}
			lc := convert(s[l])
			sFreq[lc]--
			if sFreq[lc] == tFreq[lc]-1 {
				created--
			}
			l++
		}
	}

	return ans
}

func convert76(c byte) byte {
	if c >= 'a' && c <= 'z' {
		return c - 'a'
	} else {
		return c - 'A' + 26
	}
}
