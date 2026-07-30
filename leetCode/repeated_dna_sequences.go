// https://leetcode.com/problems/repeated-dna-sequences/
// N is the length of s
// time complexity is O(N)
// space complexity is O(N)

package main

func findRepeatedDnaSequences(s string) []string {
	var ans []string
	if len(s) <= 10 {
		return ans
	}

	seen := make(map[int]int)
	charVals := map[byte]int{
		'A': 0,
		'C': 1,
		'G': 2,
		'T': 3,
	}

	var bits int
	var mask int
	for i := range 10 {
		bits <<= 2
		bits |= charVals[s[i]]
		mask <<= 2
		mask |= 0b11
	}
	seen[bits]++

	for i := 10; i < len(s); i++ {
		bits <<= 2
		bits |= charVals[s[i]]
		bits &= mask
		seen[bits]++
		if seen[bits] == 2 {
			ans = append(ans, s[i-10+1:i+1])
		}
	}

	return ans
}
