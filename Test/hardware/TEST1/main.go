package main

import (
	"fmt"
	"github.com/tarm/serial"
	"log"
	"time"
)

func main() {

	//设置串口编号
	ser := &serial.Config{Name: "COM3", Baud: 9600}

	//打开串口
	conn, err := serial.OpenPort(ser)
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan struct{})
	//启动一个协程循环发送
	go func() {
		for {
			revData := []byte("01 03 00 00 00 10 44 06")
			_, err := conn.Write(revData)
			if err != nil {
				log.Println(err)
				continue
			}
			log.Printf("Tx:%s \n", fmt.Sprintf("%s", revData))
			time.Sleep(time.Second)
		}
	}()

	//保持数据持续接收
	go func() {
		for {
			log.Println("R start\n")
			buf := make([]byte, 1024)

			lens, err := conn.Read(buf)

			if err != nil {
				log.Println("R err\n")
				log.Println(err)
				continue
			}
			revData := buf[:lens]
			log.Printf("Rx:%s \n", fmt.Sprintf("%s", revData))
			log.Println("R end\n")
		}
	}()

	<-ch


}
func f1(in chan int) {
	fmt.Println(<-in)
}
