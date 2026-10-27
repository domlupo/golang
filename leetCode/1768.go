// https://leetcode.com/problems/merge-strings-alternately/
// M is the length of word1 and N is the length of word2
// time complexity is O(M + N)
// space complexity is O(1)

package main

func mergeAlternately(word1 string, word2 string) string {
	bytes := make([]byte, len(word1)+len(word2))
	var p1, p2 int
	for p1 < len(word1) || p2 < len(word2) {
		if p1 < len(word1) {
			bytes[p1+p2] = word1[p1]
			p1++
		}
		if p2 < len(word2) {
			bytes[p1+p2] = word2[p2]
			p2++
		}
	}
	return string(bytes)
}
