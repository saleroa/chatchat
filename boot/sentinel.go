package boot

import (
	"chatchat/app/global"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.uber.org/zap"
)

func SentinelSetup() {
	err := sentinel.InitWithConfigFile("./sentinel.yaml")
	if err != nil {
		global.Logger.Fatal("initialize sentinel failed", zap.Error(err))
	}

	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "chatchat_user",
			Threshold:              200,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Throttling,
			MaxQueueingTimeMs:      1000,
			StatIntervalInMs:       10000,
		},
	})
	if err != nil {
		global.Logger.Fatal("initialize sentinel failed", zap.Error(err))
	}
	global.Logger.Info("initialize sentinel success")
}
