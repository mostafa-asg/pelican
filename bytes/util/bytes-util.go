package util

import (
	"encoding/binary"
	"fmt"
	"io"
)

func ToUint32(bytes []byte) uint32 {
	return binary.BigEndian.Uint32(bytes)
}

func ToBytes(number uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, number)
	return buf
}

func ReadBytes(numberOfBytes int, reader io.Reader) ([]byte, error) {
	buf := make([]byte, numberOfBytes)
	bytes, err := reader.Read(buf)
	if err != nil {
		return nil, err
	}
	if bytes != numberOfBytes {
		return nil, fmt.Errorf("Read %d byte(s) instead of %d byte(s)", bytes, numberOfBytes)
	}

	return buf, nil
}
