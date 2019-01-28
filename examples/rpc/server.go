package main

import (
	"github.com/yeqown/rpc"
)

type Int struct{}

// Args ... for Sum Method
type Args struct {
	A int
	B int
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
	srv := rpc.NewServerWithCodec(nil)
	srv.RegisterName(new(Int), "Add")
	srv.Start("127.0.0.1:9998", "127.0.0.1:9999")
}
