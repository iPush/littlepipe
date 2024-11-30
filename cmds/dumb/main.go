package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ipush/littlepipe/pkg/observability"
	"github.com/ipush/littlepipe/pkg/pipeline"
	pipe "github.com/ipush/littlepipe/pkg/pipeline"
	"github.com/ipush/littlepipe/pkg/sink/stdout"
	"github.com/ipush/littlepipe/pkg/source/stdin"
	"github.com/ipush/littlepipe/pkg/stage/words"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	pipeName := "dumb"
	logger := observability.NewLogger()
	tracer := observability.NewTracer(pipeName)
	metrics := observability.NewMetrics(pipeName)

	config := pipe.Config{
		BufferSize: 100,
		RetryCount: 3,
		RetryDelay: time.Second,
		Logger:     logger,
		Metrics:    metrics,
		Tracer:     tracer,
	}

	upperStage := words.NewUppercaseStage()
	obStage := pipeline.NewObservableStage(
		"uppercase",
		upperStage,
		logger,
		tracer,
		metrics,
	)
	// 创建并配置 pipeline
	pipe := pipe.NewLittlePipe(config)
	pipe.SetSource(stdin.NewStdinSource()).
		AddStage(obStage).
		SetSink(stdout.NewStdoutSink())

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()
	// 运行 pipeline
	if err := pipe.Run(); err != nil {
		log.Fatal(err)
	}
}
