package parser

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
)

type GetResult struct {
	command string
	Key     string
}

func (r *GetResult) Command() string {
	return r.command
}

type GetParser struct{}

func (p *GetParser) Accept(command []byte) bool {
	if len(command) > 4 &&
		command[0] == 'G' &&
		command[1] == 'E' &&
		command[2] == 'T' &&
		command[3] == ' ' {
		return true
	}

	return false
}

func (p *GetParser) Parse(command []byte) (ParseResult, error) {
	reader := bufio.NewReader(bytes.NewReader(command))

	line, err := reader.ReadString('\n')
	if err != nil {
		return &GetResult{}, err
	}

	tokens := strings.Split(line, " ")
	if tokens[0] != "GET" {
		return &GetResult{}, errors.New("Parse error: expected GET but found " + tokens[0])
	}

	key := strings.Trim(tokens[1], "\n")

	return &GetResult{
		command: "GET",
		Key:     key,
	}, nil
}

func NewGetParser() *GetParser {
	return &GetParser{}
}
