package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

const exchange ="go_ex"
func main()  {
	conn,err:=amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err!=nil{
		panic(err)

	}
	ch,err:=conn.Channel()

	if err!=nil{
		panic(err)
	}

	err=ch.ExchangeDeclare(
		exchange,
		 "fanout",
		 true,
		 false,
		 false		,
		false		,
		 nil,
		)

	if err!=nil{
		panic(nil)
	}
	go subscribe(conn,exchange)
	go subscribe(conn,exchange)

		var i int=0
		for  {
			i++
			err:=ch.Publish(
				exchange,
				"",
				false,
				false ,
				amqp.Publishing{
				    Body: []byte(fmt.Sprintf("msg %d\n",i)),
				})
			if err!=nil{
				fmt.Println(err.Error())
			}
			time.Sleep(time.Second)
		}

	

}
func subscribe( conn *amqp.Connection,ex string)  {
	ch,err:=conn.Channel()
	if err!=nil{
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		true,  // autoDelete
		false, // exlusive
		false, // noWait
		nil,   // args
	)
	if err!=nil{
		panic(nil)
	}
	defer ch.QueueDelete(
		q.Name,
		false,
		false,
		false,
		)
	err=ch.QueueBind(
		q.Name,
		"",
		ex,
		false,
		nil,
		)
	if err!=nil{
		fmt.Println("qname:",q.Name)
		panic(err)
	}
	consume("c1",ch,q.Name)
}
func consume(comsumer string, ch *amqp.Channel,q string)  {
	msgs,err:=ch.Consume(
		q,comsumer,
		true,
		false,
		false,
		false,
		nil,
		)
	if err!=nil{
		panic(err)
	}
	for msg :=range msgs{
		fmt.Printf("%s:%s",comsumer,msg.Body)
	}
}