// https://leetcode.com/problems/longest-common-prefix/

package main

import "sort"

func longestCommonPrefix(strs []string) string {
	sort.Strings(strs)
	first := strs[0]
	last := strs[len(strs)-1]

	var i int
	for i < len(first) && i < len(last) && first[i] == last[i] {
		i++
	}
	return first[:i]
}
