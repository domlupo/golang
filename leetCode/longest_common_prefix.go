// https://leetcode.com/problems/longest-common-prefix/

package main

func longestCommonPrefix(strs []string) string {
	for i := 0; i < len(strs)-1; i++ {
		if i >= 0 && strs[i] > strs[i+1] {
			strs[i], strs[i+1] = strs[i+1], strs[i]
			i -= 2
		}
	}

	first := strs[0]
	last := strs[len(strs)-1]
	var i int
	for i < len(first) && i < len(last) && first[i] == last[i] {
		i++
	}

	return first[:i]
}
