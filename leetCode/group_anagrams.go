// https://leetcode.com/problems/group-anagrams/
// N is length of strs
// K is max length of strs[i]
// time complexity is O(N*K)
// space complexity is O(N*K)

package main

func groupAnagrams(strs []string) [][]string {
	groups := make(map[[26]int][]string)
	var ans [][]string

	for _, str := range strs {
		var counts [26]int
		for _, c := range str {
			counts[c-'a']++
		}
		groups[counts] = append(groups[counts], str)
	}

	for _, group := range groups {
		ans = append(ans, group)
	}

	return ans
}
