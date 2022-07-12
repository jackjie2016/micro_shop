package main

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
	"log"
)
const resName = "example-flow-qps-resource"
func main()  {
	// We should initialize Sentinel first.
	//err:=sentinel.InitDefault()

	conf := config.NewDefaultConfig()
	// for testing, logging output to console
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		log.Fatal(err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               resName,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,//直接拒绝
			Threshold:              10,
			StatIntervalInMs:       1000,
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}
	for  {
		ee, eb := sentinel.Entry(resName, sentinel.WithTrafficType(base.Inbound))
		if eb!=nil{
			fmt.Println("限流了")
		}else{
			fmt.Println("正常")
			ee.Exit()
		}
	}
	//ch := make(chan struct{})
	//
	//for i := 0; i < 11; i++ {
	//	go func() {
	//		for {
	//			e, b := sentinel.Entry(resName, sentinel.WithTrafficType(base.Inbound))
	//			if b != nil {
	//				fmt.Println("限流了")
	//				// Blocked. We could get the block reason from the BlockError.
	//				time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
	//			} else {
	//				fmt.Println("正常")
	//				// Passed, wrap the logic here.
	//				time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond)
	//
	//				// Be sure the entry is exited finally.
	//				e.Exit()
	//			}
	//
	//		}
	//	}()
	//}
	//
	//// Simulate a scenario in which flow rules are updated concurrently
	//go func() {
	//	time.Sleep(time.Second * 10)
	//	fmt.Println("这是干嘛")
	//	_, err = flow.LoadRules([]*flow.Rule{
	//		{
	//			Resource:               resName,
	//			TokenCalculateStrategy: flow.Direct,
	//			ControlBehavior:        flow.Reject,
	//			Threshold:              80,
	//			StatIntervalInMs:       1000,
	//		},
	//	})
	//	if err != nil {
	//		log.Fatalf("Unexpected error: %+v", err)
	//		return
	//	}
	//}()
	//fmt.Println(<-ch)

}
