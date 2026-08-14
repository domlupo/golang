// https://leetcode.com/problems/car-fleet/
// N is the length of position and speed
// time complexity is O(N*log(N))
// space complexity is O(N)

package main

import "slices"

type drive struct {
	position, speed int
}

func carFleet(target int, position []int, speed []int) int {
	drives := make([]drive, len(position))
	for i := range len(position) {
		drives[i] = drive{position: position[i], speed: speed[i]}
	}

	slices.SortFunc(drives, func(i, j drive) int {
		return j.position - i.position
	})

	var count int
	var priorArrival float64
	for _, drive := range drives {
		arrival := float64(target-drive.position) / float64(drive.speed)
		if arrival > priorArrival {
			count++
		}
		priorArrival = max(priorArrival, arrival)
	}

	return count
}
