package transform

import (
	"github.com/expr-lang/expr"
	"github.com/ipush/littlepipe/pkg/pipeline"
)

type compiledRule struct {
	target   string
	program  *expr.Program
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
		})
	}
	return transformer
}

func (t *ExprTransform) Process(msg *pipeline.Message) (*pipeline.Message, error) {

}
