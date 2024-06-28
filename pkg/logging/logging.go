package logging

import (
	"log"
	"os"
)

var logger *log.Logger

func InitLogger(level string) {
	logger = log.New(os.Stdout, "EKOL0D: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Error(err error) {
	logger.Println("ERROR: ", err)
}

func Info(msg string) {
	logger.Println("INFO: ", msg)
}

func Warn(msg string) {
	logger.Println("WARN: ", msg)
}
