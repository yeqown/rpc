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
	// srv.Start("", "127.0.0.1:9999")
	srv.ServeTCP("127.0.0.1:9998")
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
	c := rpc.NewClientWithCodec(nil, "127.0.0.1:9998")
	testAddOverTCP(c)
	// testAddOverHTTP(c)
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
```

## Output

```sh
2019/11/30 10:29:25 [debug]: a new call
2019/11/30 10:29:25 [debug]: recv response body: ����K��*rpc.stdResponse��stdResponse��Rply
ErrErrcode��    ��
2019/11/30 10:29:25 [debug]: len(resps)=1, resp.Reply()=��
2019/11/30 10:29:25 [debug]: call done
[TCP] Int.Add(1, 222) got 223, want: 223
```