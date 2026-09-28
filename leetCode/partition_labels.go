// https://leetcode.com/problems/partition-labels/
// N is the length of s, K is the unique characters in s
// time complexity is O(N)
// space complexity is O(K)

package main

func partitionLabels(s string) []int {
	ends := make(map[rune]int)

	for i, r := range s {
		ends[r] = i
	}

	var ans []int
	var start, end int
	for i, r := range s {
		end = max(end, ends[r])
		if i == end {
			ans = append(ans, end-start+1)
			start, end = i+1, i+1
		}
	}

	return ans
}
