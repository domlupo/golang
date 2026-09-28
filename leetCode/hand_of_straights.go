// https://leetcode.com/problems/hand-of-straights/
// N is the length of hand
// time complexity is O(N)
// space complexity is O(N)

package main

func isNStraightHand(hand []int, groupSize int) bool {
	if len(hand)%groupSize != 0 {
		return false
	}

	counts := make(map[int]int)
	for _, n := range hand {
		counts[n]++
	}

	for _, card := range hand {
		start := card
		for counts[start-1] > 0 {
			start--
		}

		for start <= card {
			for counts[start] > 0 {
				for i := start; i < start+groupSize; i++ {
					if counts[i] == 0 {
						return false
					}
					counts[i]--
				}
			}
			start++
		}
	}

	return true
}
