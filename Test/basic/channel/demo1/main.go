package main

import (
	"fmt"
	"time"
)

func worker(i int, c chan int) {

	for {
		fmt.Printf("worker %d,received %d \n", i, <-c)
	}
}
func chanDemo() {
	var channels [10]chan int
	for i := 0; i < 10; i++ {
		channels[i] = make(chan int) //需要make 开辟下内存
		go worker(i, channels[i])
	}
	for i := 0; i < 10; i++ {
		channels[i] <- 'a' + i
	}

	for i := 0; i < 10; i++ {
		channels[i] <- 'A' + i
	}
	time.Sleep(time.Minute * 5)
}
func channelBuffer() {
	c := make(chan int, 3)
	go worker(0, c)
	c <- 'a'
	c <- 'b'
	c <- 'c'
	c <- 'a'
	c <- 'b'
	c <- 'c'
	time.Sleep(time.Second * 10)
}
func channelClose() {
	c := make(chan int)
	go worker(0, c)
	c <- 'a'
	c <- 'b'
	c <- 'c'
	c <- 'd'
	c <- 'e'
	close(c)
	time.Sleep(time.Second * 5)
}
func main() {
	channelBuffer()
}
