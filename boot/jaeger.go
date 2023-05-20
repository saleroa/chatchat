package boot

import (
	"chatchat/app/global"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"io"
)

func JaegerSetup() {
	tracer, closer, err := InitJaeger("chatchat")
	//tracer.StartSpan("root_tracer")
	if err != nil {
		global.Logger.Fatal("initialize jaeger failed", zap.Error(err))
	}
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	global.Logger.Info("initialize jaeger success")
}

func InitJaeger(ServiceName string) (opentracing.Tracer, io.Closer, error) {
	cfg := &config.Configuration{
		ServiceName: ServiceName,
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "127.0.0.1:6831",
		},
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, nil, err
	}
	return tracer, closer, nil
}
