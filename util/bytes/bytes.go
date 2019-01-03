package bytes

import (
	"encoding/binary"
	"fmt"
	"io"
)

func ToUint16(bytes []byte) uint16 {
	return binary.BigEndian.Uint16(bytes)
}

func ToInt16(bytes []byte) int16 {
	n := binary.BigEndian.Uint16(bytes)
	return int16(n)
}

func ToUint32(bytes []byte) uint32 {
	return binary.BigEndian.Uint32(bytes)
}

func ToInt32(bytes []byte) int32 {
	n := binary.BigEndian.Uint32(bytes)
	return int32(n)
}

func ToUint64(bytes []byte) uint64 {
	return binary.BigEndian.Uint64(bytes)
}

func ToInt64(bytes []byte) int64 {
	n := binary.BigEndian.Uint64(bytes)
	return int64(n)
}

func ToString(bytes []byte) string {
	return string(bytes)
}

func ToBool(b byte) bool {
	if b == 0 {
		return false
	}

	return true
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
