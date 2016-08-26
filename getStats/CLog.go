package getStats

import (
	"log"
)

var debug = true

func SetDebug(isDebug bool) {
	debug = isDebug
}

func Clog(format string, v ...interface{}) {
	if debug {
		log.Printf(format, v...)
	}
}