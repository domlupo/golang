// https://leetcode.com/problems/binary-watch/submissions/1909048171/

package main

import "fmt"

func readBinaryWatch(turnedOn int) []string {
	var o []string

	for h := range 12 {
		for m := range 60 {
			if countSetBits(h)+countSetBits(m) == turnedOn {
				o = append(o, fmt.Sprintf("%d:%02d", h, m))
			}
		}
	}

	return o
}

func countSetBits(n int) int {
	var count int

	for n != 0 {
		n &= (n - 1)
		count++
	}

	return count
}
