package main

import (
	"fmt"
	"log"

	"github.com/yeqown/rpc"
	"github.com/yeqown/rpc/jsonrpc"
)

// Args ...
type Args struct {
	A int `json:"a"`
	B int `json:"b"`
}

// Result .
type Result struct {
	Sum int `json:"sum"`
}

func main() {
	c := rpc.NewClientWithCodec(jsonrpc.NewJSONCodec(), "127.0.0.1:9998")

	testAddOverTCP(c)
	// testAddOverHTTP(c)
}

func testAddOverTCP(c *rpc.Client) {
	var (
		args   = &Args{A: 1, B: 222}
		result = &Result{}
	)
	if err := c.Call("Int.Add", args, &result); err != nil {
		log.Printf("c.TCP client got err: %v", err)
	}

	fmt.Printf("[TCP] Int.Add(%d, %d) got %d, want: %d\n", args.A, args.B, result.Sum, args.A+args.B)
}

// func testAddOverHTTP(c *rpc.Client) {
// 	var (
// 		args   = &Args{A: 12312, B: 8712}
// 		result = &Result{}
// 	)
// 	if err := c.HTTP("Int.Sum", args, &result); err != nil {
// 		log.Printf("c.HTTP got err: %v", err)
// 	}

// 	fmt.Printf("[HTTP] Int.Sum(%d, %d) got %d, want: %d\n", args.A, args.B, result.Sum, args.A+args.B)
// }
