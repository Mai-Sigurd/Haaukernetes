package utils

import "log"

func ErrHandler(err error) {
	if err != nil {
		log.Println(err)
	}
}

func ErrHandlerFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
