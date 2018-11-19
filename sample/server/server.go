package main

import (
	"net/http"

	"github.com/yeqown/rpc"
)

// Int ... custom type for JSON-RPC test
type Int int

// Args ... for Sum Method
type Args struct {
	A int `json:"a"`
	B int `json:"b"`
}

// Sum ...
func (i *Int) Sum(args *Args, reply *int) error {
	// println("called", args.A, args.B)
	*reply = args.A + args.B
	// *reply = 2
	return nil
}

// MultyArgs ... from Multy Int.Method
type MultyArgs struct {
	A *Args `json:"aa"`
	B *Args `json:"bb"`
}

// MultyReply ...
type MultyReply struct {
	A int `json:"aa"`
	B int `json:"bb"`
}

// Multy ... times params
func (i *Int) Multy(args *MultyArgs, reply *MultyReply) error {
	reply.A = (args.A.A * args.A.B)
	reply.B = (args.B.A * args.B.B)
	// fmt.Println(*args.A, *args.B, *reply)
	return nil
}

func main() {
	s := rpc.NewServer()
	i := new(Int)
	s.Register(i)
	go s.HandleTCP("127.0.0.1:9999")

	// to support http Request
	http.ListenAndServe(":9998", s)
}
