// https://leetcode.com/problems/last-stone-weight/

package main

func lastStoneWeight(stones []int) int {
	for len(stones) >= 2 {
		sortStones(stones)
		s1 := stones[len(stones)-1]
		s2 := stones[len(stones)-2]
		stones = stones[:len(stones)-2]
		if s1-s2 > 0 {
			stones = append(stones, s1-s2)
		}
	}

	if len(stones) == 1 {
		return stones[0]
	} else {
		return 0
	}
}

func sortStones(s []int) {
	for i := 0; i < len(s)-1; i++ {
		if i >= 0 && s[i] > s[i+1] {
			s[i], s[i+1] = s[i+1], s[i]
			i -= 2
		}
	}
}
