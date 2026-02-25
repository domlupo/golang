// https://leetcode.com/problems/find-the-town-judge/

package main

func findJudge(n int, trust [][]int) int {
	if n == 1 {
		return 1
	}

	trusts := make([]int, n+1)
	for _, t := range trust {
		trusts[t[0]]--
		trusts[t[1]]++
	}

	for ti, t := range trusts {
		if t == n-1 {
			return ti
		}
	}
	return -1
}
