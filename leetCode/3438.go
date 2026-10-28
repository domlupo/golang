// https://leetcode.com/problems/find-valid-pair-of-adjacent-digits-in-string/
// N is the length of s
// time complexity is O(N)
// space complexity is O(1)

package main

func findValidPair(s string) string {
	counts := make([]int, 10)
	for _, c := range s {
		counts[int(c-'0')]++
	}

	for i := 1; i < len(s); i++ {
		v1 := int(s[i-1] - '0')
		v2 := int(s[i] - '0')
		if s[i-1] == s[i] || v1 != counts[v1] || v2 != counts[v2] {
			continue
		}
		return s[i-1 : i+1]
	}

	return ""
}
