package parser

import (
	"bufio"
	"bytes"
	"errors"
	"strconv"
	"strings"
)

type CounterResult struct {
	command string
	Key     string
	Typ     string
	Change  int
}

func (r *CounterResult) Command() string {
	return r.command
}

type CounterParser struct{}

func (p *CounterParser) Accept(command []byte) bool {
	if len(command) > 2 &&
		command[0] == 'C' &&
		command[1] == ' ' {
		return true
	}

	return false
}

func (p *CounterParser) Parse(command []byte) (ParseResult, error) {
	reader := bufio.NewReader(bytes.NewReader(command))

	line, err := reader.ReadString('\n')
	if err != nil {
		return &CounterResult{}, err
	}

	tokens := strings.Split(line, " ")
	if tokens[0] != "C" {
		return &CounterResult{}, errors.New("Parse error: expected C but found " + tokens[0])
	}

	typ := tokens[1]
	if typ == "Inc" || typ == "Dec" {
		key := tokens[2]
		change, err := strconv.Atoi(strings.Trim(tokens[3], "\n"))
		if err != nil {
			return &CounterResult{}, err
		}

		return &CounterResult{
			command: "C",
			Key:     key,
			Typ:     typ,
			Change:  change,
		}, nil
	} else {
		return &CounterResult{}, errors.New("Parse error: invalid counter type.")
	}
}

func NewCounterParser() *CounterParser {
	return &CounterParser{}
}
