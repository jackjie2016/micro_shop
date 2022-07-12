package main

import (
	"encoding/json"
	"fmt"
)

type people struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Id   int    `json:"id"`
}

type student struct {
	people
	Id int `json:"sid"`
}

func main() {
	msg := "{\"name\":\"zhangsan\", \"age\":18, \"id\":122463, \"sid\":122464}"
	var someOne student
	var someOne2 interface{}
	if err := json.Unmarshal([]byte(msg), &someOne); err == nil {
		fmt.Println(someOne)
		fmt.Println(someOne.people)
	} else {
		fmt.Println(err)
	}

	if err := json.Unmarshal([]byte(msg), &someOne2); err == nil {
		fmt.Println(someOne2)

	} else {
		fmt.Println(err)
	}
}
