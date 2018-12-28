package parser

import (
	"bufio"
	"bytes"
	"errors"
	"strconv"
	"strings"
)

type PutResult struct {
	command string
	Key     string
	Value   []byte
}

func (r *PutResult) Command() string {
	return r.command
}

type PutParser struct{}

func (p *PutParser) Accept(command []byte) bool {
	if len(command) > 4 &&
		command[0] == 'P' &&
		command[1] == 'U' &&
		command[2] == 'T' &&
		command[3] == ' ' {
		return true
	}

	return false
}

func (p *PutParser) Parse(command []byte) (ParseResult, error) {
	reader := bufio.NewReader(bytes.NewReader(command))

	line, err := reader.ReadString('\n')
	if err != nil {
		return &PutResult{}, err
	}

	tokens := strings.Split(line, " ")
	if tokens[0] != "PUT" {
		return &PutResult{}, errors.New("Parse error: expected PUT but found " + tokens[0])
	}

	key := tokens[1]

	valueSize, err := strconv.Atoi(strings.Trim(tokens[2], "\n"))
	if err != nil {
		return &PutResult{}, err
	}

	value := make([]byte, valueSize)
	_, err = reader.Read(value)
	if err != nil {
		return &PutResult{}, err
	}

	return &PutResult{
		command: "PUT",
		Key:     key,
		Value:   value,
	}, nil
}

func NewPutParser() *PutParser {
	return &PutParser{}
}
