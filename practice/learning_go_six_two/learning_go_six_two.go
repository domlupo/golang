package main

import "fmt"

func main() {
	pets := []string{"dog", "cat"}
	fmt.Println(pets)
	updateSlice(pets, "mouse")
	fmt.Println(pets)
	growSlice(pets, "fish")
	fmt.Println(pets)
	removeAndGrowSlice(pets, "hamster")
	fmt.Println(pets)
}

func updateSlice(slice []string, s string) {
	slice[len(slice)-1] = s
	fmt.Println(slice)
}

func growSlice(slice []string, s string) {
	slice = append(slice, s)
	fmt.Println(slice)
}

func removeAndGrowSlice(slice []string, s string) {
	slice = slice[:len(slice)-1]
	slice = append(slice, s)
	fmt.Println(slice)
}
