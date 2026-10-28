// https://leetcode.com/problems/count-pairs-of-similar-strings/
// N is the length of words
// time complexity is O(N)
// space complexity is O(1) because there is at most 26 hash

package main

func similarPairs(words []string) int {
	hashs := make(map[int]int)
	for _, w := range words {
		var hash int
		for _, c := range w {
			flag := 1 << int(c-'a')
			hash |= flag
		}
		hashs[hash]++
	}

	var ans int
	for _, v := range hashs {
		for v > 1 {
			ans += (v - 1)
			v--
		}
	}
	return ans
}
