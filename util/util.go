package util

import (
	"errors"
	"log"
	"strconv"
	"time"
)

func ToTimeDuration(s string) (time.Duration, error) {
	length := len(s)
	if length < 2 {
		return 0, errors.New("At least two characters is needed")
	}

	number, err := strconv.Atoi(s[:len(s)-1])
	if err != nil {
		return 0, err
	}

	duration := s[len(s)-1:]
	switch duration {
	case "s":
		return time.Duration(number) * time.Second, nil
	case "m":
		return time.Duration(number) * time.Minute, nil
	case "h":
		return time.Duration(number) * time.Hour, nil
	default:
		return 0, errors.New("invalid time duration")
	}
}

func HandleError(err error, action func(err error)) {
	if err != nil {
		action(err)
	}
}

func Fatal(err error) {
	log.Fatal(err)
}
