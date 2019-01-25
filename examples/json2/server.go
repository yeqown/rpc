package main

import (
	"github.com/yeqown/rpc"
	"github.com/yeqown/rpc/json2"
)

type Int struct{}

// Args ... for Sum Method
type Args struct {
	A int `json:"a"`
	B int `json:"b"`
}

// Add ...
func (i *Int) Add(args *Args, reply *int) error {
	*reply = args.A + args.B
	return nil
}

// Sum ...
func (i *Int) Sum(args *Args, reply *int) error {
	*reply = args.A + args.B
	return nil
}

func main() {
	srv := rpc.NewServerWithCodec(json2.NewJSONCodec())
	srv.Register(new(Int))
	srv.Start("127.0.0.1:9998", "127.0.0.1:9999")
}
