package transform

import (
	"fmt"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/ipush/littlepipe/pkg/pipeline"
)

type compiledRule struct {
	target   string
	program  *vm.Program
	type_    pipeline.FieldType
	required bool
}

type ExprTransform struct {
	config TransformConfig
	rules  []*compiledRule
}

func NewExprTransformer(config TransformConfig) *ExprTransform {
	transformer := &ExprTransform{
		config: config,
	}

	transformer.rules = make([]*compiledRule, 0, len(config.Rules))
	for _, rule := range config.Rules {
		program, err := expr.Compile(rule.Expr,
			expr.AllowUndefinedVariables())
		if err != nil {
			// TODO: logging this error
			continue
		}

		transformer.rules = append(transformer.rules, &compiledRule{
			target:   rule.Target,
			program:  program,
			required: rule.Required,
			type_:    rule.Type,
		})
	}
	return transformer
}

func (t *ExprTransform) Process(msg *pipeline.Message) (*pipeline.Message, error) {
	newRecord := &pipeline.Record{
		Schema: &pipeline.Schema{
			Fields: make([]pipeline.Field, 0, len(t.rules)),
		},
		Data:      make(map[string]pipeline.Value),
		Timestamp: msg.Payload.Timestamp,
	}

	// prepare expr env
	env := make(map[string]any)
	for name, value := range msg.Payload.Data {
		env[name] = value.Value
	}

	for _, rule := range t.rules {
		result, err := expr.Run(rule.program, env)
		if err != nil {
			if rule.required {
				return nil, fmt.Errorf("execute rule %s: %w", rule.target, err)
			}
			continue
		}

		value, err := convertValue(result, rule.type_)
		if err != nil {
			return nil, fmt.Errorf("convert value for %s: %w", rule.target, err)
		}

		newRecord.Schema.Fields = append(newRecord.Schema.Fields, pipeline.Field{
			Name:     rule.target,
			Type:     rule.type_,
			Required: rule.required,
		})

		newRecord.Data[rule.target] = value
	}

	return &pipeline.Message{
		ID:       msg.ID,
		Payload:  newRecord,
		Metadata: msg.Metadata,
	}, nil
}

func convertValue(v interface{}, t pipeline.FieldType) (pipeline.Value, error) {
	switch t {
	case pipeline.TypeString:
		str, ok := v.(string)
		if !ok {
			return pipeline.Value{}, fmt.Errorf("cannot convert %T to string", v)
		}
		return pipeline.Value{Type: t, Value: str}, nil

	case pipeline.TypeInt64:
		switch num := v.(type) {
		case int:
			return pipeline.Value{Type: t, Value: int64(num)}, nil
		case int64:
			return pipeline.Value{Type: t, Value: num}, nil
		case float64:
			return pipeline.Value{Type: t, Value: int64(num)}, nil
		default:
			return pipeline.Value{}, fmt.Errorf("cannot convert %T to int64", v)
		}

	case pipeline.TypeFloat64:
		switch num := v.(type) {
		case float64:
			return pipeline.Value{Type: t, Value: num}, nil
		case int:
			return pipeline.Value{Type: t, Value: float64(num)}, nil
		case int64:
			return pipeline.Value{Type: t, Value: float64(num)}, nil
		default:
			return pipeline.Value{}, fmt.Errorf("cannot convert %T to float64", v)
		}

	case pipeline.TypeBoolean:
		b, ok := v.(bool)
		if !ok {
			return pipeline.Value{}, fmt.Errorf("cannot convert %T to boolean", v)
		}
		return pipeline.Value{Type: t, Value: b}, nil

	default:
		return pipeline.Value{}, fmt.Errorf("unsupported type: %v", t)
	}
}
