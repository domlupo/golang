// https://leetcode.com/problems/partition-labels/
// N is the length of s, K is the unique characters in s
// time complexity is O(N)
// space complexity is O(K)

package main

func partitionLabels(s string) []int {
	starts := make(map[byte]int)
	ends := make(map[byte]int)

	for i := 0; i < len(s); i++ {
		b := s[i]
		if _, ok := starts[b]; !ok {
			starts[b] = i
		}
		ends[b] = i
	}

	var ans []int
	var start, end int
	for i := 0; i < len(s); i++ {
		b := s[i]
		start = min(start, starts[b])
		end = max(end, ends[b])
		if i == end {
			ans = append(ans, end-start+1)
			start, end = i+1, i+1
		}
	}

	return ans
}
