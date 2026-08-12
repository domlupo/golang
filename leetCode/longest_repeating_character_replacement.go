// https://leetcode.com/problems/longest-repeating-character-replacement/
// N is the length of s
// time complexity is O(N)
// space complexity is O(26)

package main

func characterReplacement(s string, k int) int {
	var freq [26]int
	var ans, l, maxFreq int

	for r := 0; r < len(s); r++ {
		freq[s[r]-'A']++
		maxFreq = max(maxFreq, freq[s[r]-'A'])
		isValid := r+1-l-maxFreq <= k
		if !isValid {
			freq[s[l]-'A']--
			l++
		}
		ans = r + 1 - l
	}

	return ans
}
