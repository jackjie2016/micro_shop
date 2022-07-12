package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

func main()  {
	conn,err:=amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err!=nil{
		panic(err)

	}
	ch,err:=conn.Channel()

	if err!=nil{
		panic(err)
	}
	q, err:=ch.QueueDeclare(
		"go_q1",
		 false,
		 false,
		 false,
		 false		,
		 nil,
		)
	if err!=nil{
		panic(nil)
	}
	go cusume(conn,"c1",q.Name)
	go cusume(conn,"c2",q.Name)
	go cusume(conn,"c3",q.Name)
	waiting := make(chan struct{})
	go func() {
		var i int=0
		for  {
			i++
			err:=ch.Publish("",
				q.Name ,
				false,
				false ,
				amqp.Publishing{
				    Body: []byte(fmt.Sprintf("msg %d\n",i)),
				})
			if err!=nil{
				fmt.Println(err.Error())
			}
			time.Sleep(time.Second*2)
		}
	}()
	
	<-waiting
}

func cusume( conn *amqp.Connection,comsume,q string)  {
	ch,err:=conn.Channel()
	if err!=nil{
		panic(err)
	}
	msgs,err:=ch.Consume(q,comsume, true, false, false, false, nil)
	if err!=nil{
		panic(err)
	}
	for msg :=range msgs{
		fmt.Printf("%s,%s\n",comsume,msg.Body)
	}
}