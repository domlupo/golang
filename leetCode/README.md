# LeetCode

## Revisit

- 703
  - Kth Largest Element in a stream requires heap import and knowing interfaces.
- 157
  - Read N Characters Given Read4
  - This one's explanation is confusing
- 222, 408
  - This is harder than an easy
- 170, 225, 232
  - Very few companies have used this recently and it is an implement class
- 338, 645
  - Bit manipulation

## Skipped

- 175, 181, 182, 183, 196, 197, 511, 512, 577, 584, 586, 595, 596, 597, 603, 607, 610, 613, 619,
  620, 627,
  - SQL question. I think this is for data engineers?
- 193, 195
  - bash

## Tips

- Can convert a string "s" into bytes via "bytes := []byte(s)".
- Can do s = "1" + s to avoid string mutation
- To get integer do, i, _ := strconv.Atoi(string(number1[p1]))
- It is usually easiest to form strings by using s := rune[] or s := byte[] then string(s).
- string('A'+1) -> B, string('A'+2) -> C.
- strconv.Itoa(1)+"->"+strconv.Itoa(2) returns "1->2"
- s := "hello", s[5:] is "". Even though intuitively 5 is out of bounds
- odd position bit mask is 0x55555555 -> 0x55555555 & 0b11111111 -> 0b01010101
- even position bit mask:  0xAAAAAAAA -> 0xAAAAAAAA & 0b11111111 -> 0b10101010
- power of two check: (n & (n-1)) == 0. power of two means only 1 bit is set.
- When _, ok := contains[n]. ok is true if contains[n] was ever set. Even for contains[n] = false.
- fmt.Println(string('a'-32)) lower to upper
- fmt.Println(string('A'+32)) upper to lower

## Binary Search Tree

Is a tree where all values left are lower, and all values right are higher.

The structure is typically:
  call recursive function on left child of node
  do work on node
  call recursive function on right child of node

### math.MaxInt bounds

Avoid math.MaxInt when possible. math.MaxInt has weirdness near bounds for addition and subtraction.

i := -math.MaxInt
i -= 1
i *= -1 is still a negative number

### Brian Kernighan's algorithm to set count bits

func countSetBitsManual(n int) int {
	count := 0
	for n != 0 {
		n &= (n - 1) // Clears the rightmost set bit
		count++
	}
	return count
}

count = 0
111 &= 110 -> 110
count = 1
110 &= 101 -> 100
count = 2
100 &= 011 -> 000
count = 3
