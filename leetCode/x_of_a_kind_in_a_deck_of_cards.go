// https://leetcode.com/problems/x-of-a-kind-in-a-deck-of-cards/

package main

func hasGroupsSizeX(deck []int) bool {
	counts := make(map[int]int)

	for _, n := range deck {
		counts[n]++
	}

	var g int
	for _, c := range counts {
		g = gcd(g, c)
	}

	return g > 1
}

func gcd(v1, v2 int) int {
	for v2 != 0 {
		v1, v2 = v2, v1%v2
	}
	return v1
}
