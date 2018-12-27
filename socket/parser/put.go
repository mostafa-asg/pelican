package parser

import (
	"bufio"
	"bytes"
	"errors"
	"strconv"
	"strings"
)

// ParsePut parses
func ParsePut(data []byte) (string, []byte, error) {
	reader := bufio.NewReader(bytes.NewReader(data))

	line, err := reader.ReadString('\n')
	if err != nil {
		return "", nil, err
	}

	tokens := strings.Split(line, " ")
	if tokens[0] != "PUT" {
		return "", nil, errors.New("Parse error: expected PUT but found " + tokens[0])
	}

	key := tokens[1]

	valueSize, err := strconv.Atoi(strings.Trim(tokens[2], "\n"))
	if err != nil {
		return "", nil, err
	}

	value := make([]byte, valueSize)
	n, err := reader.Read(value)
	if err != nil {
		return "", nil, err
	}

	return key, value[:n], nil
}
