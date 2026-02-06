// https://leetcode.com/problems/assign-cookies/

package main

func findContentChildren(g []int, s []int) int {
	for i := 1; i < len(g); i++ {
		if i >= 1 && g[i] > g[i-1] {
			g[i], g[i-1] = g[i-1], g[i]
			i -= 2
		}
	}

	for i := 1; i < len(s); i++ {
		if i >= 1 && s[i] > s[i-1] {
			s[i], s[i-1] = s[i-1], s[i]
			i -= 2
		}
	}

	var count int
	var p int
	for _, greed := range g {
		if p > len(s)-1 {
			break
		}
		if s[p] >= greed {
			count++
			p++
		}
	}
	return count
}
