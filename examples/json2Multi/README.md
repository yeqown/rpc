# JSON RPC 2.0 Multi

## Server-side

```golang
// example/rpc/server.go
package main

import (
	"github.com/yeqown/rpc"
	"github.com/yeqown/rpc/json2"
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
	srv := rpc.NewServerWithCodec(json2.NewJSONCodec())
	srv.Register(new(Int))
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
	"github.com/yeqown/rpc/json2"
)

// Args ...
type Args struct {
	A int `json:"a"`
	B int `json:"b"`
}

func main() {
	c := rpc.NewClientWithCodec(json2.NewJSONCodec(), "127.0.0.1:9998", "127.0.0.1:9999")
	testAddOverTCP(c)
	testAddOverHTTP(c)
}

func testAddOverTCP(c *rpc.Client) {
	var (
		sum  int
		args = &Args{A: 1, B: 222}
	)
	if err := c.CallOverTCP("Int.Add", args, &sum); err != nil {
		println("got err: ", err.Error())
	}

	fmt.Printf("[TCP] Int.Add(%d, %d) got %d, want: %d\n", args.A, args.B, sum, args.A+args.B)
}

func testAddOverHTTP(c *rpc.Client) {
	var (
		sum  int
		args = &Args{A: 12312, B: 8712}
	)
	if err := c.CallOverHTTP("Int.Sum", args, &sum); err != nil {
		println("got err: ", err.Error())
	}

	fmt.Printf("[HTTP] Int.Sum(%d, %d) got %d, want: %d\n", args.A, args.B, sum, args.A+args.B)
}
```

## Output

```sh
[TCP] Int.Add(1, 222) got 223, want: 223
2019/01/25 16:26:08 [debug]: send request [addr: 127.0.0.1:9999] [data: eyJpZCI6ImQxMTcxZjMzNGQ3NmFlN2JkYjBlNDI2M2NkMDM4NWNjIiwibWV0aG9kIjoiSW50LlN1bSIsInBhcmFtcyI6IlpYbEthRWxxYjNoTmFrMTRUV2wzYVZscFNUWlBSR040VFc0d1BRPT0iLCJqc29ucnBjIjoiMi4wIn0=]
2019/01/25 16:26:08 [debug]: got response eyJpZCI6ImQxMTcxZjMzNGQ3NmFlN2JkYjBlNDI2M2NkMDM4NWNjIiwicmVzdWx0IjoiWlhsS2VtUlhNR2xQYWtsNFRVUkpNR1pSUFQwPSIsImpzb25ycGMiOiIyLjAifQ==
2019/01/25 16:26:08 [debug]: client decode to reply: origin eyJzdW0iOjIxMDI0fQ==, decoded: 0xc00009a5d8
[HTTP] Int.Sum(12312, 8712) got 21024, want: 21024
```