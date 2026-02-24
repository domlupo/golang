// https://leetcode.com/problems/number-of-recent-calls/

package main

type RecentCounter struct {
	pings []int
}

func (this *RecentCounter) Ping(t int) int {
	this.pings = append(this.pings, t)

	for i := 0; i < len(this.pings); i++ {
		if this.pings[i] >= t-3000 {
			this.pings = this.pings[i:]
			break
		}
	}

	return len(this.pings)
}
