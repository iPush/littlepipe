package main

import (
	"log"
	"time"

	pipe "github.com/ipush/littlepipe"
	"github.com/ipush/littlepipe/sink/stdout"
	"github.com/ipush/littlepipe/source/stdin"
	"github.com/ipush/littlepipe/stage/words"
)

func main() {
	config := pipe.Config{
		BufferSize: 100,
		RetryCount: 3,
		RetryDelay: time.Second,
	}

	// 创建并配置 pipeline
	pipe := pipe.NewLittlePipe(config)
	pipe.SetSource(stdin.NewStdinSource()).
		AddStage(words.NewUppercaseStage()).
		AddStage(words.NewFilterEmptyStage()).
		SetSink(stdout.NewStdoutSink())

	// 运行 pipeline
	if err := pipe.Run(); err != nil {
		log.Fatal(err)
	}
}
