package main

import (
	"fmt"
	"math/rand"
	"time"
)

func msgGen(name string) chan string {
	c := make(chan string)
	go func() {
		i := 0
		for {
			time.Sleep(time.Duration(rand.Intn(2000)) * time.Millisecond)
			c <- fmt.Sprintf("service name :%s message %d", name, i)
			i++
		}
	}()
	return c
}

func fanIn(c1, c2 chan string) chan string {
	c := make(chan string)
	go func() {
		for {
			select {
			case m := <-c1:
				c <- m
			case m := <-c2:
				c <- m
			}
		}

	}()

	return c
}

func fanIns(chs ...chan string) chan string {
	c := make(chan string)
	i := 1
	for _, ch := range chs {
		go func(in chan string) {
			for {
				c <- <-in
				i = i + 1
				fmt.Println(i)
			}

		}(ch)
	}

	return c
}
func main() {
	m1 := msgGen("1")
	m2 := msgGen("2")
	m3 := msgGen("3")
	m := fanIns(m1, m2, m3)
	for {
		fmt.Println(<-m)
	}
}
