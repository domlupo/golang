// https://leetcode.com/problems/construct-the-rectangle/

package main

import "math"

func constructRectangle(area int) []int {
	var l int
	w := int(math.Sqrt(float64(area))) + 1
	for w > 0 {
		w -= 1
		if area%w == 0 {
			l = area / w
			break
		}
	}
	return []int{l, w}
}
