// https://leetcode.com/problems/most-common-word/

package main

import "math"

func mostCommonWord(paragraph string, banned []string) string {
	banLookUp := make(map[string]bool)
	for _, w := range banned {
		banLookUp[w] = true
	}

	nonBanCount := make(map[string]int)
	maxCount := math.MinInt
	var mostCommon string
	for i := 0; i < len(paragraph); i++ {
		b := []byte{}
		for i < len(paragraph) && validChar(paragraph[i]) {
			b = append(b, lower1(paragraph[i]))
			i += 1
		}

		if len(b) == 0 {
			continue
		}

		w := string(b)

		if banLookUp[w] {
			continue
		}

		nonBanCount[w]++
		if nonBanCount[w] > maxCount {
			maxCount = nonBanCount[w]
			mostCommon = w
		}
	}

	return mostCommon
}

func validChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')
}

func lower1(b byte) byte {
	if b >= 'a' && b <= 'z' {
		return b
	} else {
		return b + 32
	}
}
