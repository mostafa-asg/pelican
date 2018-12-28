package parser

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
)

type DelResult struct {
	command string
	Key     string
}

func (r *DelResult) Command() string {
	return r.command
}

type DelParser struct{}

func (p *DelParser) Accept(command []byte) bool {
	if len(command) > 4 &&
		command[0] == 'D' &&
		command[1] == 'E' &&
		command[2] == 'L' &&
		command[3] == ' ' {
		return true
	}

	return false
}

func (p *DelParser) Parse(command []byte) (ParseResult, error) {
	reader := bufio.NewReader(bytes.NewReader(command))

	line, err := reader.ReadString('\n')
	if err != nil {
		return &DelResult{}, err
	}

	tokens := strings.Split(line, " ")
	if tokens[0] != "DEL" {
		return &DelResult{}, errors.New("Parse error: expected DEL but found " + tokens[0])
	}

	key := strings.Trim(tokens[1], "\n")

	return &DelResult{
		command: "DEL",
		Key:     key,
	}, nil
}

func NewDelParser() *DelParser {
	return &DelParser{}
}
