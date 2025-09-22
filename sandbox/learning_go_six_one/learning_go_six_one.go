package main

import "fmt"

type Person struct {
	firstName string
	lastName  string
	age       int
}

func main() {
	p1 := makePerson("dom1", "lupo1", 1)
	p2 := makePersonPointer("dom2", "lupo2", 2)
	fmt.Println(p1)
	fmt.Println(p2)
}

func makePerson(firstName string, lastName string, age int) Person {
	p3 := Person{
		firstName,
		lastName,
		age,
	}
	return p3
}

func makePersonPointer(firstName string, lastName string, age int) *Person {
	p4 := Person{
		firstName,
		lastName,
		age,
	}
	return &p4
}
