package main

import (
	"fmt"
	"rpc"
)

func main() {
	c := rpc.NewClient()
	c.DialTCP("127.0.0.1:9999")

	rcv := c.Send("hhh")
	fmt.Println("Recived from server: ", rcv)
}
