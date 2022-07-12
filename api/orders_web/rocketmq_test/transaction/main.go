package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"google.golang.org/grpc/codes"
	"os"
	"strconv"
	"time"
)
type OrderListener struct {
	Code codes.Code
	Detail string

}

/*
订单的实际处理逻辑写在这边
 */
func (o *OrderListener)ExecuteLocalTransaction(*primitive.Message) primitive.LocalTransactionState{
	fmt.Println("本地事务开始")
	time.Sleep(time.Second*3)
	fmt.Printf("本地事务成功")
	return primitive.UnknowState
}

// When no response to prepare(half) message. broker will send check message to check the transaction status, and this
// method will be invoked to get local transaction status.
func (o *OrderListener)CheckLocalTransaction(*primitive.MessageExt) primitive.LocalTransactionState{
	fmt.Println("消息会查")
	return primitive.RollbackMessageState
}
func main() {
	p, err := rocketmq.NewTransactionProducer(&OrderListener{},
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"43.242.33.9:9876"})),
		producer.WithRetry(1),
	)

	if err!=nil{
		panic(err.Error())
	}
   err = p.Start()

	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	topic := "order_reback"

	for i := 0; i < 10; i++ {
		msg := &primitive.Message{
			Topic: topic,
			Body:  []byte("事务消息 " + strconv.Itoa(i)),
		}
		res, err := p.SendMessageInTransaction(context.Background(), msg)

		if err != nil {
			fmt.Printf("发送失败: %s\n", err)
		} else {
			fmt.Printf("发送成功: result=%s\n", res.String())
		}
	}
	err = p.Shutdown()
	if err != nil {
		fmt.Printf("shutdown producer error: %s", err.Error())
	}
}