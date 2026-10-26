// https://leetcode.com/problems/longest-common-prefix/
// S is the total characters in the array of strings
// time complexity is O(S)
// space complexity is O(1)

package main

func longestCommonPrefix(strs []string) string {
	var i int
	for i < len(strs[0]) {
		for j := range strs {
			s := strs[j]
			if i >= len(s) || strs[0][i] != s[i] {
				return strs[0][:i]
			}
		}
		i++
	}
	return strs[0]
}
