package parser

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
)

// ParseDel parses
func ParseDel(data []byte) (string, error) {
	reader := bufio.NewReader(bytes.NewReader(data))

	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	tokens := strings.Split(line, " ")
	if tokens[0] != "DEL" {
		return "", errors.New("Parse error: expected DEL but found " + tokens[0])
	}

	key := strings.Trim(tokens[1], "\n")

	return key, nil
}
