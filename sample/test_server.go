package main

import (
	// "fmt"
	"net/http"
	"rpc"
)

type Int int

type Args struct {
	A int `json:"a"`
	B int `json:"b"`
}

func (i *Int) Sum(args *Args, reply *int) error {
	// println("called", args.A, args.B)
	*reply = args.A + args.B
	// *reply = 2
	return nil
}

type MultyArgs struct {
	A *Args `json:"aa"`
	B *Args `json:"bb"`
}

type MultyReply struct {
	A int `json:"aa"`
	B int `json:"bb"`
}

func (i *Int) Multy(args *MultyArgs, reply *MultyReply) error {
	reply.A = (args.A.A * args.A.B)
	reply.B = (args.B.A * args.B.B)
	// fmt.Println(*args.A, *args.B, *reply)
	return nil
}

func main() {
	s := rpc.NewServer()
	mine_int := new(Int)
	s.Register(mine_int)
	go s.HandleTCP("127.0.0.1:9999")

	// to support http Request
	http.ListenAndServe(":9998", s)
}
