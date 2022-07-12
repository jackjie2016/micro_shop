package main

import "fmt"

func main()  {
	ch := make(chan struct{})
	go func() {
		for  {
			fmt.Println("22")
		}
	}()
	<-ch
}
