package utils

import (
	"log"
	"os"
	"time"
)

func SetLog() {
	currentTime := time.Now()
	f, err := os.OpenFile("logs/"+currentTime.Format("2006.01.02 15:04:05"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
}
