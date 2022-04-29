package tracer

import (
	"io"
	"time"
	"web_app/settings"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

var (
	Tracer opentracing.Tracer
)

func InitTracer(cfg *settings.JaegerConfig) error {
	jaegerTracer, _, err := NewJaegerTracer(cfg.ServiceName, cfg.AgentHostPort)
	if err != nil {
		return err
	}
	Tracer = jaegerTracer
	return nil
}


func NewJaegerTracer(serviceName, agentHostPort string) (opentracing.Tracer, io.Closer, error){
	cfg := &config.Configuration{
		ServiceName:         serviceName,
		Sampler:             &config.SamplerConfig{
			Type: "const",
			Param: 1, // 全采样
		},
		Reporter:            &config.ReporterConfig{
			LogSpans: true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort: agentHostPort,
		},
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, nil, err
	}
	opentracing.SetGlobalTracer(tracer) // 设置全局的 tracer对象
	return tracer, closer, nil
}