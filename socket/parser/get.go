package parser

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
)

// ParseGet parses
func ParseGet(data []byte) (string, error) {
	reader := bufio.NewReader(bytes.NewReader(data))

	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	tokens := strings.Split(line, " ")
	if tokens[0] != "GET" {
		return "", errors.New("Parse error: expected GET but found " + tokens[0])
	}

	key := strings.Trim(tokens[1], "\n")

	return key, nil
}
