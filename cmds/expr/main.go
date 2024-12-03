package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ipush/littlepipe/pkg/pipeline"
	"github.com/ipush/littlepipe/pkg/stage/transform"
)

func main() {
	// 创建转换器配置
	config := transform.TransformConfig{
		Rules: []transform.TransformRule{
			{
				Target:   "full_name",
				Expr:     `first_name + " " + last_name`,
				Type:     pipeline.TypeString,
				Required: true,
			},
			{
				Target:   "age",
				Expr:     "2024 - birth_year",
				Type:     pipeline.TypeInt64,
				Required: true,
			},
		},
	}

	// 创建转换器
	transformer := transform.NewExprTransformer(config)

	// 创建输入数据
	inputRecord := &pipeline.Record{
		Schema: &pipeline.Schema{
			Fields: []pipeline.Field{
				{Name: "first_name", Type: pipeline.TypeString},
				{Name: "last_name", Type: pipeline.TypeString},
				{Name: "birth_year", Type: pipeline.TypeInt64},
			},
		},
		Data: map[string]pipeline.Value{
			"first_name": {Type: pipeline.TypeString, Value: "John"},
			"last_name":  {Type: pipeline.TypeString, Value: "Doe"},
			"birth_year": {Type: pipeline.TypeInt64, Value: int64(1990)},
		},
		Timestamp: time.Now(),
	}

	// 创建输入消息
	inputMsg := pipeline.NewMessage(inputRecord)

	// 处理消息
	outputMsg, err := transformer.Process(inputMsg)
	if err != nil {
		log.Fatal(err)
	}

	// 输出结果
	fmt.Printf("Full Name: %v\n", outputMsg.Payload.Data["full_name"].Value)
	fmt.Printf("Age: %v\n", outputMsg.Payload.Data["age"].Value)
}
