package utils

import (
	"log"
	"runtime"

	"github.com/laxamore/mineralos/utils/Log"
)

func CheckErr(err error) {
	if err != nil {
		log.SetFlags(0)
		_, filename, lineno, ok := runtime.Caller(1)
		if ok {
			Log.Panicf("%v:%v: %v", filename, lineno, err)
		} else {
			Log.Panicf("%v", err)
		}
	}
}
