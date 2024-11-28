package stdout

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ipush/littlepipe/pkg/pipeline"
)

// StdoutSink 将数据写入标准输出
type StdoutSink struct {
	writer *bufio.Writer
}

func NewStdoutSink() *StdoutSink {
	return &StdoutSink{
		writer: bufio.NewWriter(os.Stdout),
	}
}

func (s *StdoutSink) Write(data *pipeline.Message) error {
	if str, ok := data.Payload.(string); ok {
		_, err := s.writer.WriteString(str + "\n")
		if err != nil {
			return err
		}
		return s.writer.Flush()
	}
	return fmt.Errorf("StdoutSink: unsupported data type: %T", data)
}
