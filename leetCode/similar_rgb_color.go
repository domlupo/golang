// https://leetcode.com/problems/similar-rgb-color/

package main

import "strconv"

func similarRGB(color string) string {
	shorthandVals := map[string]int{
		"00": 0,
		"11": 16*1 + 1,
		"22": 16*2 + 2,
		"33": 16*3 + 3,
		"44": 16*4 + 4,
		"55": 16*5 + 5,
		"66": 16*6 + 6,
		"77": 16*7 + 7,
		"88": 16*8 + 8,
		"99": 16*9 + 9,
		"aa": 16*10 + 10,
		"bb": 16*11 + 11,
		"cc": 16*12 + 12,
		"dd": 16*13 + 13,
		"ee": 16*14 + 14,
		"ff": 16*15 + 15,
	}

	first := hexToVal(color[1])*16 + hexToVal(color[2])
	second := hexToVal(color[3])*16 + hexToVal(color[4])
	third := hexToVal(color[5])*16 + hexToVal(color[6])

	firstColor, secondColor, thirdColor := "00", "00", "00"
	for k, v := range shorthandVals {
		if absDiff(v, first) < absDiff(shorthandVals[firstColor], first) {
			firstColor = k
		}
		if absDiff(v, second) < absDiff(shorthandVals[secondColor], second) {
			secondColor = k
		}
		if absDiff(v, third) < absDiff(shorthandVals[thirdColor], third) {
			thirdColor = k
		}
	}

	return "#" + firstColor + secondColor + thirdColor
}

func absDiff(x, y int) int {
	if x > y {
		return x - y
	} else {
		return y - x
	}
}

func hexToVal(c byte) int {
	v, err := strconv.Atoi(string(c))
	if err == nil {
		return v
	} else {
		return int(10 + byte(c) - byte('a'))
	}
}
