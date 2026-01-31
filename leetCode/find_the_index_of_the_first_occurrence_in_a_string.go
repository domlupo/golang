// https://leetcode.com/problems/find-the-index-of-the-first-occurrence-in-a-string/

package main

func strStr(haystack string, needle string) int {
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:len(needle)+i] == needle {
			return i
		}
	}
	return -1
}
