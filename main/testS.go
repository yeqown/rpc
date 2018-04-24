package main

import (
	"rpc"
)

func main() {
	s := rpc.NewServer()
	s.HandleTCP("127.0.0.1:9999")
}
