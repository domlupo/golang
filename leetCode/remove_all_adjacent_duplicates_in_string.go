// https://leetcode.com/problems/remove-all-adjacent-duplicates-in-string/

package main

func removeDuplicates1(s string) string {
	var ans []byte
	for i := 0; i < len(s); i++ {
		if len(ans) > 0 && ans[len(ans)-1] == s[i] {
			ans = ans[:len(ans)-1]
		} else {
			ans = append(ans, s[i])
		}
	}
	return string(ans)
}
