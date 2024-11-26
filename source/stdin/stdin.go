package stdin

import (
	"bufio"
	"io"
	"os"
)

type StdinSource struct {
	scanner *bufio.Scanner
}

func NewStdinSource() *StdinSource {
	return &StdinSource{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (s *StdinSource) Read() (interface{}, error) {
	if s.scanner.Scan() {
		return s.scanner.Text(), nil
	}
	if err := s.scanner.Err(); err != nil {
		return nil, err
	}
	return nil, io.EOF
}
