package main

import (
	"fmt"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"
	"log"
	"math/rand"
	"time"
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
			TokenCalculateStrategy: flow.WarmUp,
			ControlBehavior:        flow.Reject,//直接拒绝
			Threshold:              100,//并发数
			WarmUpPeriodSec: 30,       //30s达到1000的并发数
		},
	})
	if err != nil {
		log.Fatalf("Unexpected error: %+v", err)
		return
	}
     var globalTotal,PassTotal,blockTotal int
	ch := make(chan struct{})

	for i := 0; i < 11; i++ {
		go func() {
			for {
				globalTotal++
				e, b := sentinel.Entry(resName, sentinel.WithTrafficType(base.Inbound))
				if b != nil {
					blockTotal++
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond*5)
				} else {
					PassTotal++
					time.Sleep(time.Duration(rand.Uint64()%10) * time.Millisecond*5)
					e.Exit()
				}

			}
		}()
	}

	// Simulate a scenario in which flow rules are updated concurrently

	go func() {
		var OldPass,OldTotal,OldBlock int  //每秒的情况
		for  {
			sectotal:=globalTotal-OldTotal
			OldTotal=globalTotal

			secPass:=PassTotal-OldPass
			OldPass=PassTotal


			secBlock:=blockTotal-OldBlock
			OldBlock=blockTotal


			time.Sleep(time.Second)
			fmt.Printf("total:%d,passtotal:%d,blocktotal:%d \n",sectotal,secPass,secBlock)
		}

	}()
	<-ch


}
