package words

import (
	"fmt"
	"strings"

	"github.com/ipush/littlepipe/pkg/pipeline"
)

// UppercaseStage 将文本转换为大写
type UppercaseStage struct{}

func NewUppercaseStage() *UppercaseStage {
	return &UppercaseStage{}
}

func (s *UppercaseStage) Process(data *pipeline.Message) (*pipeline.Message, error) {
	if str, ok := data.Payload.(string); ok {
		return pipeline.NewMessage(strings.ToUpper(str)), nil
	}
	return nil, fmt.Errorf("UppercaseStage: unsupported data type: %T", data)
}
