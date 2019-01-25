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
	var (
		args   = &Args{A: 1, B: 222}
		result = &Result{}
	)
	if err := c.CallOverTCP("Int.Add", args, &result); err != nil {
		log.Printf("c.CallOverTCP client got err: %v", err)
	}

	fmt.Printf("[TCP] Int.Add(%d, %d) got %d, want: %d\n", args.A, args.B, result.Sum, args.A+args.B)
}

func testAddOverHTTP(c *rpc.Client) {
	var (
		args   = &Args{A: 12312, B: 8712}
		result = &Result{}
	)
	if err := c.CallOverHTTP("Int.Sum", args, &result); err != nil {
		log.Printf("c.CallOverHTTP got err: %v", err)
	}

	fmt.Printf("[HTTP] Int.Sum(%d, %d) got %d, want: %d\n", args.A, args.B, result.Sum, args.A+args.B)
}
