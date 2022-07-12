package main

import (
	"context"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
)

func main() {

	p, err := rocketmq.NewProducer(
		producer.WithNsResolver(primitive.NewPassthroughResolver([]string{"47.97.215.11:9876"})),
		producer.WithRetry(2),
	)

	if err != nil {
		zap.S().Errorf("生成 producer 失败: %s", err.Error())
		return
	}
	var msg *primitive.Message

	var orderInfo OrderInfo
	orderInfo.OrderAmount = 12
	orderInfo.ID = 1
	_ = json.Unmarshal(msg.Body, &orderInfo)

	msg2 := &primitive.Message{
		Topic: "order_timeout",
		Body:  msg.Body,
	}

	msg2.WithDelayTimeLevel(4) //跟普通比就多一句这个
	_, err = p.SendSync(context.Background(), msg2)
	if err != nil {
		zap.S().Errorf("发送信息失败: %s", err.Error())
		return
	}
}

type OrderListener struct {
	Code        codes.Code
	Detail      string
	ID          int32
	OrderAmount float32
	ctx         context.Context
}
type OrderInfo struct {
	ID          int32
	OrderAmount float32
}
