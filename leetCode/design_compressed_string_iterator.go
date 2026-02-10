// https://leetcode.com/problems/design-compressed-string-iterator/

package main

import "strconv"

type StringIterator struct {
	str   string
	ptr   int
	ch    byte
	count int
}

func Constructor2(compressedString string) StringIterator {
	s := StringIterator{
		str: compressedString,
	}
	return s
}

func (this *StringIterator) Next() byte {
	if this.count == 0 && this.ptr < len(this.str) {
		this.ch = this.str[this.ptr]
		this.ptr++
		start := this.ptr
		for this.ptr < len(this.str) && this.str[this.ptr] >= '0' && this.str[this.ptr] <= '9' {
			this.ptr++
		}
		count, _ := strconv.Atoi(this.str[start:this.ptr])
		this.count = count
	}

	if this.count > 0 {
		this.count--
		return this.ch
	}

	return ' '
}

func (this *StringIterator) HasNext() bool {
	return this.count != 0 || this.ptr < len(this.str)
}
