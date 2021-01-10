package server

import (
	"flag"
)

// ParseArgs - parse api arguments
func ParseArgs() int {
	var port int
	flag.IntVar(&port, "port", 8090, "Specify server port")
	flag.Parse()
	return port
}