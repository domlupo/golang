// https://leetcode.com/problems/rectangle-overlap/

package main

func isRectangleOverlap(rec1 []int, rec2 []int) bool {
	x1 := rec1[0]
	x2 := rec1[2]
	x3 := rec2[0]
	x4 := rec2[2]

	y1 := rec1[1]
	y2 := rec1[3]
	y3 := rec2[1]
	y4 := rec2[3]

	innerLeft := max(x1, x3)
	innerRight := min(x2, x4)
	innerDown := max(y1, y3)
	innerUp := min(y2, y4)

	return innerLeft < innerRight && innerDown < innerUp
}
