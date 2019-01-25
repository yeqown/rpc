package main

import (
	"fmt"

	"github.com/yeqown/rpc"
)

// Args ...
type Args struct {
	A int
	B int
}

func main() {
	c := rpc.NewClientWithCodec(nil, "127.0.0.1:9998", "127.0.0.1:9999")
	testAddOverTCP(c)
	testAddOverHTTP(c)
}

func testAddOverTCP(c *rpc.Client) {
	var (
		sum  int
		args = &Args{A: 1, B: 222}
	)
	if err := c.Call("Int.Add", args, &sum); err != nil {
		println("got err: ", err.Error())
	}

	fmt.Printf("[TCP] Int.Add(%d, %d) got %d, want: %d\n", args.A, args.B, sum, args.A+args.B)
}

func testAddOverHTTP(c *rpc.Client) {
	var (
		sum  int
		args = &Args{A: 1111, B: 222}
	)
	if err := c.CallHTTP("Int.Add", args, &sum); err != nil {
		println("got err: ", err.Error())
	}

	fmt.Printf("[HTTP] Int.Add(%d, %d) got %d, want: %d\n", args.A, args.B, sum, args.A+args.B)
}
