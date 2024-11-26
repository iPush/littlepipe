package stdout

import (
	"bufio"
	"fmt"
	"os"
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

func (s *StdoutSink) Write(data interface{}) error {
	if str, ok := data.(string); ok {
		_, err := s.writer.WriteString(str + "\n")
		if err != nil {
			return err
		}
		return s.writer.Flush()
	}
	return fmt.Errorf("StdoutSink: unsupported data type: %T", data)
}
