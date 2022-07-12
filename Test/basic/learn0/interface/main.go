package main

import (
	"fmt"
	"strconv"
)

type Binary uint64

func (i Binary) String() string {
	return strconv.FormatUint(uint64(i), 10)
}

type Stringer interface {
	String() string
}

func test(s Stringer) {

	fmt.Println(s.String())
}

func main() {
	b := Binary(0x123)
	test(b)

}
