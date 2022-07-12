package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"math/cmplx"
)

func variableZeroValue() {
	var a int
	var b string
	fmt.Printf("%d,%s\n", a, b)
}

//欧拉公示
func euler() {
	//e 的 1i*Π次
	//1i 是虚数
	fmt.Println(
		cmplx.Pow(math.E, 1i*math.Pi) + 1)
	fmt.Println(fmt.Sprintf("%.3f", cmplx.Pow(math.E, 1i*math.Pi)+1))
	c := 3 + 4i
	fmt.Println(cmplx.Abs(c))

}

func triangle() {
	var a, b int = 3, 4
	var c int
	c = int(math.Sqrt(float64(math.Pow(float64(a), 2) + math.Pow(float64(b), 2))))
	fmt.Println(c)
}

//常量可以不指定类型
func constf() {
	const (
		a, b = 3, 4
	)

	var c int
	c = int(math.Sqrt(float64(math.Pow(a, 2) + math.Pow(b, 2))))
	fmt.Println(c)
}

//自增 iota的使用
func enums() {
	const (
		ppp = 1
		cpp = iota
		java
		golang
		php
		python
	)

	fmt.Println(ppp, cpp, java, golang, php, python)

	const (
		b = 1 << (iota * 10)
		kb
		mb
		gb
		tb
		pb
	)
	fmt.Println(b, kb, mb, gb, tb, pb)
}

func readFile() {
	const filename = "abc.txt"
	if contents, err := ioutil.ReadFile(filename); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s\n", contents)
	}
}

func grade(scroe int) string {
	g := ""
	switch {
	case scroe <= 100:
		g = "A"
		fallthrough
	case scroe < 90:
		g = "B"
	case scroe < 80:
		g = "C"
	case scroe < 70:
		g = "F"
	case scroe < 60:
		g = "F"

	default:
		panic(fmt.Sprintf("wrong score %s\n", scroe))

	}
	return g
}
func main() {
	//variableZeroValue()
	//euler()
	//triangle()
	//constf()
	//enums()
	//readFile()
	//fmt.Println(grade(90))

	size := 2
	page := 2
	params := make([]int, 0)
	params = append(params, size, page)
	fmt.Println(params)
}
