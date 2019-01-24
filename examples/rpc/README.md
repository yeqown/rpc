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
	srv := rpc.NewServerWithCodec("127.0.0.1:9999", nil)
    srv.RegisterName(new(Int), "Add")
    // srv.Register(new(Int)) will register all exported methods.
	srv.ServeTCP()
}

```

## Client-side

```golang
// examples/rpc/client.go
package main

import (
	"fmt"
	"reflect"

	"github.com/yeqown/rpc"
)

// Args ...
type Args struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {
	c := rpc.NewClient("127.0.0.1:9999")
	testAdd(c)
}

func testAdd(c *rpc.Client) {
	var (
		sum  int
		args = &Args{A: 1, B: 222}
	)
	if err := c.Call("Int.Add", args, &sum); err != nil {
		println("got err: ", err.Error())
	}

	fmt.Printf("Int.Add(%d, %d) got %d, want: %d\n", args.A, args.B, sum, args.A+args.B)
}
```

## Output

```sh
Int.Add(1, 222) got 223, want: 223
```