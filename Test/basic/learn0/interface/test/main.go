package main

import "fmt"

type Person struct {
	Age int
}

func (p Person) GetNum() int{
	return p.Age
}
func Add(s ...interface{}) int {
	total:=0
	for _,v:=range s{
		person := v.(Person)
		total+= person.GetNum()
	}
	return total
}

func main()  {
	//不能使用	Persons:=[]Person{} 涉及interface底层原理
	Persons:=[]interface{}{}

	Persons=append(Persons,Person{18})
	Persons=append(Persons,Person{19})
	fmt.Println(Add(Persons...))
}
