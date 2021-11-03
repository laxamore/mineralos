package Log

import (
	"log"

	"github.com/gin-gonic/gin"
)

func Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func Print(format string) {
	log.Print(format)
}

func Panic(format string, v ...interface{}) {
	if gin.Mode() == gin.TestMode {
		log.Print(format)
	} else {
		log.Panic(format)
	}
}

func Panicf(format string, v ...interface{}) {
	if gin.Mode() == gin.TestMode {
		log.Printf(format, v...)
	} else {
		log.Panicf(format, v...)
	}
}
