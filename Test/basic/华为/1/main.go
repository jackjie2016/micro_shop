package main

import (
	"fmt"
	"math"
)

var k int64
var count int

func input() {
	fmt.Scanln(&k)
	//fmt.Println(2*math.Pow(10, 9) + 14)
	if k < 1 || k > (2*int64(math.Pow(10, 9))+14) {
		input()
	}
}

func main() {
	count = 0
	input()

	if !do(k) {
		fmt.Print(k)
	}
	// 	fmt.Println(slince[0:count])
}

func do(number int64) bool {
	fmt.Printf("number:%d, count：%d\n", number, count)
	var i int64
	flag := false
	stop := int64(math.Sqrt(float64(number)))
	fmt.Printf("number:%d, count：%d,stop:%d\n", number, count, stop)
	for i = 1; i < stop+1; i = i + 2 {
		var j int64
		if i == 1 {
			j = 2
		} else {
			j = i
		}
		if number%j == 0 {

			fmt.Print(j)
			fmt.Print(" ")
			flag = true
			//fmt.Printf(" count：%d,slince len %d,cap %d,v:%v\n", count, len(slince), cap(slince), slince)
			if !do(int64(number / j)) {
				fmt.Print(number / j)
			}

			break
			//fmt.Printf(" count：%d,slince len %d,cap %d,v:%v\n", count, len(slince), cap(slince), slince)
		}
	}
	return flag

}
