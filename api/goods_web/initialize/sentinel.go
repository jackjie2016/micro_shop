package initialize
import (
	"go.uber.org/zap"
	"log"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"

)
func InitSentinel()  {
	conf := config.NewDefaultConfig()
	// for testing, logging output to console
	conf.Sentinel.Log.Logger = logging.NewConsoleLogger()
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		log.Fatal(err)
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "goods-list",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,//直接拒绝
			Threshold:              5,
			StatIntervalInMs:       5000,
		},
		{
			Resource:               "goods-update",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,//直接拒绝
			Threshold:              5,
			StatIntervalInMs:       5000,
		},
		{
			Resource:               "goods-add",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,//直接拒绝
			Threshold:              5,
			StatIntervalInMs:       5000,
		},
		{
			Resource:               "goods-detail",
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,//直接拒绝
			Threshold:              5,
			StatIntervalInMs:       5000,
		},
	})
	if err != nil {
		zap.S().Fatalf("Unexpected error: %+v", err)
		return
	}
}