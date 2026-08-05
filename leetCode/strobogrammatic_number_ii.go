// https://leetcode.com/problems/strobogrammatic-number-ii/
// N is the passed n
// time and space complexity is hard to calculate but it is close to below.
// time complexity is O(N * 5^N)
// space complexity is O(N * 5^N)

package main

import "strconv"

func findStrobogrammatic(n int) []string {
	candidates := genCandidates(n)
	var ans []string
	for _, candidate := range candidates {
		n, _ := strconv.Atoi(candidate)
		s := strconv.Itoa(n)
		if len(s) == len(candidate) {
			ans = append(ans, candidate)
		}
	}
	return ans
}

func genCandidates(n int) []string {
	if n == 0 {
		return []string{""}
	} else if n == 1 {
		return []string{"0", "1", "8"}
	}
	var candidates []string
	for _, n := range genCandidates(n - 2) {
		candidates = append(candidates, "0"+n+"0")
		candidates = append(candidates, "1"+n+"1")
		candidates = append(candidates, "6"+n+"9")
		candidates = append(candidates, "8"+n+"8")
		candidates = append(candidates, "9"+n+"6")
	}
	return candidates
}
