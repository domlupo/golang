package main

import (
	"fmt"
	"math/rand/v2"
)

func main() {
	problems := []int{1, 9, 13, 14, 20, 21, 26, 27, 28, 35, 58, 66, 67, 69, 70, 83}

	i := rand.IntN(len(problems))
	fmt.Printf("Random problem: %d\n", problems[i])
}
