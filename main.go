package main

import "github.com/SleepWalker/sinx/snet"

func main() {
	server := snet.NewServer("Test")
	server.Serve()
}
