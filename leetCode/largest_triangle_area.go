// https://leetcode.com/problems/largest-triangle-area/

package main

import "math"

func largestTriangleArea(points [][]int) float64 {
	maxArea := float64(-math.MaxInt)
	for i := 0; i < len(points)-2; i++ {
		for j := i + 1; j < len(points)-1; j++ {
			for k := j + 1; k < len(points); k++ {
				p1x := points[i][0]
				p2x := points[j][0]
				p3x := points[k][0]
				p1y := points[i][1]
				p2y := points[j][1]
				p3y := points[k][1]

				area := 0.5 * float64(p1x*(p2y-p3y)+
					p2x*(p3y-p1y)+p3x*(p1y-p2y))
				if area < 0 {
					area *= -1
				}

				maxArea = max(maxArea, area)
			}
		}
	}

	return maxArea
}
