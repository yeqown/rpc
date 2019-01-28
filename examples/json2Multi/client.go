package main

import (
	"fmt"
	"log"

	"github.com/yeqown/rpc"
	"github.com/yeqown/rpc/json2"
)

// Args ...
type Args struct {
	A int `json:"a"`
	B int `json:"b"`
}

type Result struct {
	Sum int `json:"sum"`
}

func main() {
	c := rpc.NewClientWithCodec(json2.NewJSONCodec(), "127.0.0.1:9998", "127.0.0.1:9999")
	// c := rpc.NewClientWithCodec(json2.NewStdJSONCodec(), "127.0.0.1:9998", "127.0.0.1:9999")
	testAddOverTCP(c)
	testAddOverHTTP(c)
}

func testAddOverTCP(c *rpc.Client) {
	cfgs := []*rpc.RequestConfig{
		&rpc.RequestConfig{
			Method: "Int.Add",
			Args:   &Args{121233, 1912109},
			Reply:  &Result{},
		},
		&rpc.RequestConfig{
			Method: "Int.Sum",
			Args:   &Args{2311231, 1909},
			Reply:  &Result{},
		},
	}

	if err := c.CallOverTCPMulti(cfgs); err != nil {
		log.Printf("c.CallOverTCPMulti client got err: %v", err)
	}
	for _, cfg := range cfgs {
		args := cfg.Args.(*Args)
		result := cfg.Reply.(*Result)
		fmt.Printf("[TCP] Int.Add(%d, %d) got %d, want: %d\n", args.A, args.B, result.Sum, args.A+args.B)
	}
}

func testAddOverHTTP(c *rpc.Client) {
	cfgs := []*rpc.RequestConfig{
		&rpc.RequestConfig{
			Method: "Int.Add",
			Args:   &Args{10, 1909},
			Reply:  &Result{},
		},
		&rpc.RequestConfig{
			Method: "Int.Sum",
			Args:   &Args{21312, 1909},
			Reply:  &Result{},
		},
	}
	if err := c.CallOverHTTPMulti(cfgs); err != nil {
		log.Printf("c.CallOverHTTPMulti client got err: %v", err)
	}
	for _, cfg := range cfgs {
		args := cfg.Args.(*Args)
		result := cfg.Reply.(*Result)
		fmt.Printf("[HTTP] Int.Add(%d, %d) got %d, want: %d\n", args.A, args.B, result.Sum, args.A+args.B)
	}
}
