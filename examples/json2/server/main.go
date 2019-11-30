package main

import (
	"github.com/yeqown/rpc"
	"github.com/yeqown/rpc/jsonrpc"
)

// Int .
type Int struct{}

// Args ... for Sum Method
type Args struct {
	A int `json:"a"`
	B int `json:"b"`
}

// Result .
type Result struct {
	Sum int `json:"sum"`
}

// Add ...
func (i *Int) Add(args *Args, reply *Result) error {
	reply.Sum = args.A + args.B
	return nil
}

// Sum ...
func (i *Int) Sum(args *Args, reply *Result) error {
	reply.Sum = args.A + args.B
	return nil
}

func main() {
	srv := rpc.NewServerWithCodec(jsonrpc.NewJSONCodec())
	// srv := rpc.NewServerWithCodec(json2.NewStdJSONCodec())
	srv.Register(new(Int))
	go srv.ServeTCP("127.0.0.1:9998")
	srv.ListenAndServe("127.0.0.1:9999")
}
