package util

import "log"

func HandleError(err error, action func(err error)) {
	if err != nil {
		action(err)
	}
}

func Fatal(err error) {
	log.Fatal(err)
}
