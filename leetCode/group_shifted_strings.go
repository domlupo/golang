// https://leetcode.com/problems/group-shifted-strings/
// N is the length of strings, K is the length of the longest string
// time complexity is O(N * K)
// space complexity is O(N * K)

package main

func groupStrings(strings []string) [][]string {
	shifts := make(map[string]int)
	var ans [][]string
	for _, str := range strings {
		shift := rune(str[0] - 'a')
		var shiftedChars []rune
		for _, c := range str {
			shiftedC := c - shift
			if shiftedC < 'a' {
				shiftedC += 26
			}
			shiftedChars = append(shiftedChars, shiftedC)
		}
		shiftedStr := string(shiftedChars)
		if i, ok := shifts[shiftedStr]; ok {
			ans[i] = append(ans[i], str)
		} else {
			ans = append(ans, []string{str})
			shifts[shiftedStr] = len(ans) - 1
		}
	}
	return ans
}
