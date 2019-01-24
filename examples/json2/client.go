package main

import (
	"fmt"
	"reflect"

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
		sum     = 0
		wantSum = 3
	)
	if err := c.Call("Int.Add", &Args{A: 1, B: 2}, &sum); err != nil {
		println("got err: ", err.Error())
	}

	if !reflect.DeepEqual(sum, wantSum) {
		println(fmt.Sprintf("Int.Add Result %d not equal to %d", sum, wantSum))
		return
	}
	println("testAdd passed")
}
