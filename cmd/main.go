package main

import (
	"github.com/vlkorniienko/api/server"
)

func main() {
	port := server.ParseArgs()

	api := server.New(port)
	api.Start()
	var a chan bool
	<-a
}
