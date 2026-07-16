// https://leetcode.com/problems/longest-common-prefix/
// N is the length of strs
// time complexity is O(N*log(N))
// space complexity is O(log(N)) from slices.Sort
// TODO: Revisit there are better solutions

package main

import "slices"

func longestCommonPrefix(strs []string) string {
	slices.Sort(strs)

	first := strs[0]
	last := strs[len(strs)-1]
	var i int
	for i < len(first) && i < len(last) && first[i] == last[i] {
		i++
	}

	return first[:i]
}
