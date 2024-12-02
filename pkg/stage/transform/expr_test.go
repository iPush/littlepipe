package transform

import (
	"fmt"
	"testing"
	"time"

	"github.com/ipush/littlepipe/pkg/pipeline"
)

// createTestMessage 创建测试消息
func createTestMessage() *pipeline.Message {
	record := &pipeline.Record{
		Schema: &pipeline.Schema{
			Fields: []pipeline.Field{
				{Name: "first_name", Type: pipeline.TypeString},
				{Name: "last_name", Type: pipeline.TypeString},
				{Name: "age", Type: pipeline.TypeInt64},
				{Name: "salary", Type: pipeline.TypeFloat64},
				{Name: "is_active", Type: pipeline.TypeBoolean},
			},
		},
		Data: map[string]pipeline.Value{
			"first_name": {Type: pipeline.TypeString, Value: "John"},
			"last_name":  {Type: pipeline.TypeString, Value: "Doe"},
			"age":        {Type: pipeline.TypeInt64, Value: int64(30)},
			"salary":     {Type: pipeline.TypeFloat64, Value: float64(50000)},
			"is_active":  {Type: pipeline.TypeBoolean, Value: true},
		},
		Timestamp: time.Now(),
	}
	return pipeline.NewMessage(record)
}

// createTestTransformer 创建测试转换器
func createTestTransformer() pipeline.Stage {
	config := TransformConfig{
		Rules: []TransformRule{
			{
				Target:   "full_name",
				Expr:     `first_name + " " + last_name`,
				Type:     pipeline.TypeString,
				Required: true,
			},
			{
				Target:   "salary_after_tax",
				Expr:     "salary * 0.8",
				Type:     pipeline.TypeFloat64,
				Required: true,
			},
			{
				Target:   "can_retire",
				Expr:     "age >= 60",
				Type:     pipeline.TypeBoolean,
				Required: true,
			},
			{
				Target:   "status",
				Expr:     `if is_active { "active" } else { "inactive" }`,
				Type:     pipeline.TypeString,
				Required: true,
			},
			{
				Target:   "yearly_review",
				Expr:     `if salary >= 50000 { "excellent" } else if salary >= 30000 { "good" } else { "normal" }`,
				Type:     pipeline.TypeString,
				Required: true,
			},
		},
	}
	return NewExprTransformer(config)
}

func BenchmarkExprTransformer(b *testing.B) {
	transformer := createTestTransformer()
	msg := createTestMessage()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := transformer.Process(msg)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// 测试不同复杂度的表达式
func BenchmarkExprTransformer_Complexity(b *testing.B) {
	testCases := []struct {
		name string
		rule TransformRule
	}{
		{
			name: "Simple_Concat",
			rule: TransformRule{
				Target: "result",
				Expr:   `first_name + " " + last_name`,
				Type:   pipeline.TypeString,
			},
		},
		{
			name: "Simple_Math",
			rule: TransformRule{
				Target: "result",
				Expr:   "salary * 0.8",
				Type:   pipeline.TypeFloat64,
			},
		},
		{
			name: "Simple_Condition",
			rule: TransformRule{
				Target: "result",
				Expr:   "age >= 30",
				Type:   pipeline.TypeBoolean,
			},
		},
		{
			name: "Complex_Condition",
			rule: TransformRule{
				Target: "result",
				Expr:   `if age >= 30 && salary >= 50000 { "senior" } else if age >= 25 { "mid" } else { "junior" }`,
				Type:   pipeline.TypeString,
			},
		},
		{
			name: "Complex_Math",
			rule: TransformRule{
				Target: "result",
				Expr:   "((salary * 0.8) + (age * 1000)) / 2",
				Type:   pipeline.TypeFloat64,
			},
		},
	}

	msg := createTestMessage()

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			config := TransformConfig{
				Rules: []TransformRule{tc.rule},
			}
			transformer := NewExprTransformer(config)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := transformer.Process(msg)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// 测试不同数量的规则
func BenchmarkExprTransformer_RuleCount(b *testing.B) {
	testCases := []struct {
		name      string
		ruleCount int
	}{
		{"Rules_1", 1},
		{"Rules_5", 5},
		{"Rules_10", 10},
		{"Rules_20", 20},
		{"Rules_50", 50},
	}

	baseRule := TransformRule{
		Target: "result",
		Expr:   `first_name + " " + last_name`,
		Type:   pipeline.TypeString,
	}

	msg := createTestMessage()

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			rules := make([]TransformRule, tc.ruleCount)
			for i := 0; i < tc.ruleCount; i++ {
				rules[i] = TransformRule{
					Target: fmt.Sprintf("result_%d", i),
					Expr:   baseRule.Expr,
					Type:   baseRule.Type,
				}
			}

			transformer := NewExprTransformer(TransformConfig{Rules: rules})

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := transformer.Process(msg)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// 测试并发性能
func BenchmarkExprTransformer_Parallel(b *testing.B) {
	transformer := createTestTransformer()
	msg := createTestMessage()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := transformer.Process(msg)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
