package main

import "fmt"

func Print1to20() interface{} {
	return 200
}

func main()  {

	var str2 string = "01 03 00 00 00 10 44 06"

	var data []byte = []byte(str2)
	fmt.Println(data)

	revData := []byte("01 03 00 00 00 10 44 06")
	fmt.Println(revData[0])
	fmt.Println(revData[:])
	str:= string(revData)
	fmt.Println(str)
}