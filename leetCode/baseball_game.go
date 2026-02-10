// https://leetcode.com/problems/baseball-game/

package main

import "strconv"

func calPoints(operations []string) int {
	var record []int

	for _, op := range operations {
		switch op {
		case "+":
			record = append(record, record[len(record)-1]+record[len(record)-2])
		case "D":
			record = append(record, record[len(record)-1]*2)
		case "C":
			record = record[:len(record)-1]
		default:
			r, _ := strconv.Atoi(op)
			record = append(record, r)
		}
	}

	var sum int
	for _, r := range record {
		sum += r
	}
	return sum
}
