package main

import (
	"rpc"
)

type Args struct {
	A int `json:"a"`
	B int `json:"b"`
}

type MultyArgs struct {
	A *Args `json:"aa"`
	B *Args `json:"bb"`
}

type MultyReply struct {
	A int `json:"aa"`
	B int `json:"bb"`
}

func main() {
	c := rpc.NewClient()
	c.DialTCP("127.0.0.1:9999")

	var sum int
	c.Call("1", "Int.Sum", &Args{A: 1, B: 2}, &sum)
	println(sum)

	c.DialTCP("127.0.0.1:9999")
	var reply MultyReply
	c.Call("2", "Int.Multy", &MultyArgs{A: &Args{1, 2}, B: &Args{3, 4}}, &reply)
	println(reply.A, reply.B)
}
