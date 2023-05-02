package utils

import (
	"log"
	"os"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	TestLogger  *log.Logger
)

func SetLog() {
	currentTime := time.Now()

	fileName := "logs/" + currentTime.Format("15:04:05")

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	InfoLogger = log.New(f, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(f, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	log.SetOutput(f)
}

func SetLogTest(fileName string, normalLogging bool) {
	f, err := os.OpenFile("logs/"+fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	TestLogger = log.New(f, "", log.Ldate|log.Ltime|log.Lshortfile)
	if normalLogging {
		SetLog()
	}
}

func ErrLogger(err error) {
	if err != nil {
		ErrorLogger.Println(err)
	}
}
