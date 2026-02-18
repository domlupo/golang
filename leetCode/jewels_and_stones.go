// https://leetcode.com/problems/jewels-and-stones/

package main

func numJewelsInStones(jewels string, stones string) int {
	isJewel := make(map[rune]bool)
	for _, j := range jewels {
		isJewel[j] = true
	}

	var count int
	for _, s := range stones {
		if isJewel[s] {
			count++
		}
	}
	return count
}
