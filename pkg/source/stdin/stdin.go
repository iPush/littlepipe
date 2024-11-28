package stdin

import (
	"bufio"
	"io"
	"os"

	"github.com/ipush/littlepipe/pkg/pipeline"
)

type StdinSource struct {
	scanner *bufio.Scanner
}

func NewStdinSource() *StdinSource {
	return &StdinSource{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (s *StdinSource) Read() (*pipeline.Message, error) {
	if s.scanner.Scan() {
		return pipeline.NewMessage(s.scanner.Text()), nil
	}
	if err := s.scanner.Err(); err != nil {
		return nil, err
	}
	return nil, io.EOF
}
