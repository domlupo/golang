// https://leetcode.com/problems/moving-average-from-data-stream/

package main

type MovingAverage struct {
	stream []int
	size   int
	sum    int
}

func Constructor1(size int) MovingAverage {
	return MovingAverage{
		size: size,
	}
}

func (this *MovingAverage) Next(val int) float64 {
	this.stream = append(this.stream, val)
	this.sum += val
	if len(this.stream) > this.size {
		this.sum -= this.stream[0]
		this.stream = this.stream[1:]
	}
	return float64(this.sum) / float64(len(this.stream))
}
