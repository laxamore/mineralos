package logger

import (
	"log"
)

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Print(format ...interface{}) {
	log.Print(format...)
}

func Error(format ...interface{}) {
	log.Print(format...)
}

func Errorf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Panic(format ...interface{}) {
	log.Panic(format...)
}

func Panicf(format string, v ...interface{}) {
	log.Panicf(format, v...)
}
