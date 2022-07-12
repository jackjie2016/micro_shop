package main

import "fmt"

var z *int

func main() {
	a := 1
	z = &a
	fmt.Println(*z)
}
