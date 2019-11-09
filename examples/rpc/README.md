# RPC example

## Server-side

```golang
// example/rpc/server.go
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
```

## Client-side

```golang
// examples/rpc/client.go
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
	if err := c.TCP("Int.Add", args, &sum); err != nil {
		println("got err: ", err.Error())
	}

	fmt.Printf("[TCP] Int.Add(%d, %d) got %d, want: %d\n", args.A, args.B, sum, args.A+args.B)
}

func testAddOverHTTP(c *rpc.Client) {
	var (
		sum  int
		args = &Args{A: 1111, B: 222}
	)
	if err := c.HTTP("Int.Add", args, &sum); err != nil {
		println("got err: ", err.Error())
	}

	fmt.Printf("[HTTP] Int.Add(%d, %d) got %d, want: %d\n", args.A, args.B, sum, args.A+args.B)
}

```

## Output

```sh
[TCP] Int.Add(1, 222) got 223, want: 223
2019/01/25 11:17:56 [debug]: send request [addr: 127.0.0.1:9999] [data: Lv+DAwEBDmRlZmF1bHRSZXF1ZXN0Af+EAAECAQRNdGhkAQwAAQRBcmdzAQoAAABK/4QBB0ludC5BZGQBPEh2K0JBd0VCQkVGeVozTUIvNElBQVFJQkFVRUJCQUFCQVVJQkJBQUFBQXYvZ2dIK0NLNEIvZ0c4QUE9PQA=]
2019/01/25 11:17:56 [debug]: got response Ov+BAwEBD2RlZmF1bHRSZXNwb25zZQH/ggABAwEEUnBseQEKAAEDRXJyAQwAAQdFcnJjb2RlAQQAAAAN/4IBCEJRUUEvZ3BxAA==
[HTTP] Int.Add(1111, 222) got 1333, want: 1333
```