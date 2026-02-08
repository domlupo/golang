// https://leetcode.com/problems/distribute-candies/

package main

func distributeCandies(candyType []int) int {
	var uniqCount int
	uniq := make(map[int]bool)

	for _, c := range candyType {
		_, ok := uniq[c]
		if !ok {
			uniqCount++
			uniq[c] = true
		}
	}

	return min(uniqCount, len(candyType)/2)
}
