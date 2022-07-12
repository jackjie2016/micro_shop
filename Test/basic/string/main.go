package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	var a int
	a = 10
	a, _ = strconv.Atoi("a")
	b := strconv.Itoa(10)
	c, _ := strconv.ParseInt("AF", 16, 8)
	d := strconv.FormatInt(10, 16)
	fmt.Println(a, b, c, d)

	//string 标准库
	s := "hello world this hero"
	f1 := strings.Contains("hello world this hero", "this")
	f2 := strings.ContainsAny("hello world this hero", "twh")
	f3 := strings.ContainsRune("hello world this hero", 32)
	by := []byte(s)
	s1 := string(by[0])
	fmt.Println(f1, f2, f3, by, s1)
}
