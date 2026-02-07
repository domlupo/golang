// https://leetcode.com/problems/relative-ranks/

package main

import (
	"math"
	"strconv"
)

func findRelativeRanks(score []int) []string {
	m := -math.MaxInt
	for _, s := range score {
		m = max(m, s)
	}

	scoreToIndex := make([]int, m+1)
	for si, s := range score {
		// add 1 to differentiate score 0 and Golang array zero value of 0
		scoreToIndex[s] = si + 1
	}

	o := make([]string, len(score))
	count := 1
	for i := len(scoreToIndex) - 1; i >= 0; i-- {
		if scoreToIndex[i] != 0 {
			var str string
			switch count {
			case 1:
				str = "Gold Medal"
			case 2:
				str = "Silver Medal"
			case 3:
				str = "Bronze Medal"
			default:
				str = strconv.Itoa(count)
			}
			o[scoreToIndex[i]-1] = str
			count++
		}
	}

	return o
}
