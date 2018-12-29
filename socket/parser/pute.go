package parser

import (
	"bufio"
	"bytes"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/mostafa-asg/pelican/store"
	"github.com/mostafa-asg/pelican/util"
)

type PutEResult struct {
	command    string
	Key        string
	Value      []byte
	Expiration time.Duration
	Strategy   store.Strategy
}

func (r *PutEResult) Command() string {
	return r.command
}

type PutEParser struct{}

func (p *PutEParser) Accept(command []byte) bool {
	if len(command) > 5 &&
		command[0] == 'P' &&
		command[1] == 'U' &&
		command[2] == 'T' &&
		command[3] == 'E' &&
		command[4] == ' ' {
		return true
	}

	return false
}

func (p *PutEParser) Parse(command []byte) (ParseResult, error) {
	reader := bufio.NewReader(bytes.NewReader(command))

	line, err := reader.ReadString('\n')
	if err != nil {
		return &PutEResult{}, err
	}

	tokens := strings.Split(line, " ")
	if len(tokens) != 5 {
		return &PutEResult{}, errors.New("Not enough information in the header")
	}
	if tokens[0] != "PUTE" {
		return &PutEResult{}, errors.New("Parse error: expected PUT but found " + tokens[0])
	}

	key := tokens[1]

	expiration, err := util.ToTimeDuration(tokens[2])
	if err != nil {
		return &PutEResult{}, err
	}

	strategy, err := strconv.Atoi(tokens[3])
	if err != nil {
		return &PutEResult{}, err
	}

	valueSize, err := strconv.Atoi(strings.Trim(tokens[4], "\n"))
	if err != nil {
		return &PutEResult{}, err
	}

	value := make([]byte, valueSize)
	_, err = reader.Read(value)
	if err != nil {
		return &PutEResult{}, err
	}

	return &PutEResult{
		command:    "PUTE",
		Key:        key,
		Value:      value,
		Expiration: expiration,
		Strategy:   store.Strategy(strategy),
	}, nil
}

func NewPutEParser() *PutEParser {
	return &PutEParser{}
}
