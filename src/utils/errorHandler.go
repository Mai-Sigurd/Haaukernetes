package utils

import "log"

func ErrHandlerFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
