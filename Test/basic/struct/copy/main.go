package main

import "fmt"

type A1 struct {
	Name string
	Age  int
}
type B1 struct {
	Name string
	Age  int
}

func main() {
	a := A1{
		Name: "222",
		Age:  1,
	}
	var b B1

	b = B1(a)

	fmt.Println(b)

}
