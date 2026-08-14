// https://leetcode.com/problems/min-stack/
// Push, Pop, Top and GetMin functions are all O(1)

package main

type MinStack struct {
	values    []int
	mins      []int
	minCounts []int
}

func Constructor155() MinStack {
	return MinStack{}
}

func (this *MinStack) Push(value int) {
	this.values = append(this.values, value)

	if len(this.mins) == 0 || value < this.mins[len(this.mins)-1] {
		this.mins = append(this.mins, value)
		this.minCounts = append(this.minCounts, 1)
	} else if this.mins[len(this.mins)-1] == value {
		this.minCounts[len(this.minCounts)-1]++
	}
}

func (this *MinStack) Pop() {
	v := this.values[len(this.values)-1]
	this.values = this.values[:len(this.values)-1]

	if v == this.mins[len(this.mins)-1] {
		this.minCounts[len(this.minCounts)-1]--
	}

	if this.minCounts[len(this.minCounts)-1] == 0 {
		this.mins = this.mins[:len(this.mins)-1]
		this.minCounts = this.minCounts[:len(this.minCounts)-1]
	}
}

func (this *MinStack) Top() int {
	return this.values[len(this.values)-1]
}

func (this *MinStack) GetMin() int {
	return this.mins[len(this.mins)-1]
}
