package utils

import (
	"github.com/laxamore/mineralos/internal/logger"
	"log"
	"runtime"
)

func CheckErr(err error) {
	if err != nil {
		log.SetFlags(0)
		_, filename, lineno, ok := runtime.Caller(1)
		if ok {
			logger.Panicf("%v:%v: %v", filename, lineno, err)
		} else {
			logger.Panicf("%v", err)
		}
	}
}
