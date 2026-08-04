// https://leetcode.com/problems/rectangle-area/
// time complexity is O(1)
// space complexity is O(1)

package main

func computeArea(ax1 int, ay1 int, ax2 int, ay2 int, bx1 int, by1 int, bx2 int, by2 int) int {
	lClamp := max(ax1, bx1)
	dClamp := max(ay1, by1)
	uClamp := min(ay2, by2)
	rClamp := min(ax2, bx2)
	var overlap int
	if rClamp > lClamp && uClamp > dClamp {
		overlap = (rClamp - lClamp) * (uClamp - dClamp)
	}
	first := (ax2 - ax1) * (ay2 - ay1)
	second := (bx2 - bx1) * (by2 - by1)
	return first + second - overlap
}
