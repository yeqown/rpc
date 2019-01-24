package main

import (
	"fmt"

	"github.com/yeqown/rpc"
	"github.com/yeqown/rpc/json2"
)

// Args ...
type Args struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {
	c := rpc.NewClientWithCodec("127.0.0.1:9999", json2.NewJSONCodec())
	testAdd(c)
}

func testAdd(c *rpc.Client) {
	var (
		sum  int
		args = &Args{A: 1, B: 222}
	)
	if err := c.Call("Int.Sum", args, &sum); err != nil {
		println(`c.Call("Int.Sum", args, &sum) got err: `, err.Error())
	}

	fmt.Printf("Int.Add(%d, %d) got %d, want: %d\n", args.A, args.B, sum, args.A+args.B)
}
