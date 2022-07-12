package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/logging"
	"github.com/alibaba/sentinel-golang/util"
)

type stateChangeTestListener struct {
}

func (s *stateChangeTestListener) OnTransformToClosed(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.steategy: %+v, From %s to Closed, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToOpen(prev circuitbreaker.State, rule circuitbreaker.Rule, snapshot interface{}) {
	fmt.Printf("rule.steategy: %+v, From %s to Open, snapshot: %.2f, time: %d\n", rule.Strategy, prev.String(), snapshot, util.CurrentTimeMillis())
}

func (s *stateChangeTestListener) OnTransformToHalfOpen(prev circuitbreaker.State, rule circuitbreaker.Rule) {
	fmt.Printf("rule.steategy: %+v, From %s to Half-Open, time: %d\n", rule.Strategy, prev.String(), util.CurrentTimeMillis())
}

func main() {
	conf := config.NewDefaultConfig()
	// for testing, logging output to console
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		log.Fatal(err)
	}
	ch := make(chan struct{})
	// Register a state change listener so that we could observer the state change of the internal circuit breaker.
	circuitbreaker.RegisterStateChangeListeners(&stateChangeTestListener{})

	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		// Statistic time span=5s, recoveryTimeout=3s, maxErrorRatio=40%
		{
			Resource:                     "abc",
			Strategy:                     circuitbreaker.ErrorRatio,
			RetryTimeoutMs:               3000,
			MinRequestAmount:             10,
			StatIntervalMs:               5000,
			//StatSlidingWindowBucketCount: 10,
			Threshold:                    0.4,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	var total,block,pass ,totalerr int
	logging.Info("[CircuitBreaker ErrorRatio] Sentinel Go circuit breaking demo is running. You may see the pass/block metric in the metric log.")
	go func() {
		for {
			total++
			e, b := sentinel.Entry("abc")
			if b != nil {
				block++
				// g1 blocked
				time.Sleep(time.Duration(rand.Uint64()%20) * time.Millisecond)
			} else {
				pass++
				if rand.Uint64()%20 > 6 {
					totalerr++
					// Record current invocation as error.
					sentinel.TraceError(e, errors.New("biz error"))
				}
				// g1 passed
				time.Sleep(time.Duration(rand.Uint64()%80+20) * time.Millisecond)
				e.Exit()
			}
		}
	}()
	go func() {
		for {
			total++
			e, b := sentinel.Entry("abc")
			if b != nil {
				block++
				// g2 blocked
				time.Sleep(time.Duration(rand.Uint64()%20) * time.Millisecond)
			} else {
				pass++
				// g2 passed
				time.Sleep(time.Duration(rand.Uint64()%80+40) * time.Millisecond)
				e.Exit()
			}
		}
	}()
	go func() {
		for  {
			time.Sleep(time.Second)
			fmt.Println(float64(totalerr)/float64(total))
		}
	}()
	<-ch
}