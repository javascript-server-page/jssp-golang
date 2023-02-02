package main

import (
	"jssp/server"
)

func main() {
	jssp := new(server.JsspServer)
	jssp.Init()
	jssp.Run()
}
