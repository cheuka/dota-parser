package getStats

import "log"

var debug = true

func SetDebug(isDebug bool)  {
	debug = isDebug
}

func Clog(format string, i ...interface{}){
	if debug{
		log.Printf(format, i)
	}
}