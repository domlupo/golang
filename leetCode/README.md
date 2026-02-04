# LeetCode

## Revisit

- 125
  - Valid Palindrome
  - This one is complex and requires remembering standard library functions.
- 157
  - Read N Characters Given Read4
  - This one's explanation is confusing
- 222
  - This is harder than an easy
- 170, 225, 232
  - Very few companies have used this recently and it is an implement class
- 338
  - Bit manipulation

## Skipped

- 175, 181, 182, 183, 196, 197
  - SQL question. I think this is for data engineers?
- 193, 195
  - bash

## Tips

- It is usually easiest to form strings by using r := rune[] then string(r)
- string(rune('A'+1)) -> B, string(rune('A'+2)) -> C.
- strconv.Itoa(1)+"->"+strconv.Itoa(2) returns "1->2"
- s := "hello", s[5:] == "". Even though intuitively 5 is out of bounds
