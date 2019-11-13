# JSON RPC 2.0

About [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification)

[MultiRequest example](../json2-array)

## Server-side

```golang
package main

import (
	"github.com/yeqown/rpc"
	"github.com/yeqown/rpc/jsonrpc"
)

// Int .
type Int struct{}

// Args ... for Sum Method
type Args struct {
	A int `json:"a"`
	B int `json:"b"`
}

// Result .
type Result struct {
	Sum int `json:"sum"`
}

// Add ...
func (i *Int) Add(args *Args, reply *Result) error {
	reply.Sum = args.A + args.B
	return nil
}

// Sum ...
func (i *Int) Sum(args *Args, reply *Result) error {
	reply.Sum = args.A + args.B
	return nil
}

func main() {
	srv := rpc.NewServerWithCodec(jsonrpc.NewJSONCodec())
	// srv := rpc.NewServerWithCodec(json2.NewStdJSONCodec())
	srv.Register(new(Int))
	go srv.ServeTCP("127.0.0.1:9998")
	srv.ListenAndServe("127.0.0.1:9999")
}

```

## Client-side

```golang
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

```

## Output

`server`
```sh
2019/11/13 16:55:35 [debug]: RPC server over TCP is listening: 127.0.0.1:9998
2019/11/13 16:55:42 [debug]: recv a new request: [12 255 131 2 1 2 255 132 0 1 16 0 0 62 255 132 0 1 15 42 114 112 99 46 115 116 100 82 101 113 117 101 115 116 255 133 3 1 1 10 115 116 100 82 101 113 117 101 115 116 1 255 134 0 1 2 1 4 77 116 104 100 1 12 0 1 4 65 114 103 115 1 10 0 0 0 56 255 134 53 1 7 73 110 116 46 65 100 100 1 41 30 255 129 3 1 1 4 65 114 103 115 1 255 130 0 1 2 1 1 65 1 4 0 1 1 66 1 4 0 0 0 9 255 130 1 2 1 254 1 188 0 0]
2019/11/13 16:55:42 [debug]: [service.call] argv: &{1 222}, replyv: 0xc0001361f8
2019/11/13 16:55:42 [debug]: server called end
2019/11/13 16:55:42 [debug]: s.call(req) req: [0xc00014a360] result: [0xc00014a840]
2019/11/13 16:55:42 [debug]: ReadTCP error: EOF
```

`client`
```sh
2019/11/13 16:55:42 [debug]: a new call
2019/11/13 16:55:42 [debug]: recv response body:
                                                 ����K��*rpc.stdResponse��
    stdResponse��Rply
Err
   Errcode
          ��	��
2019/11/13 16:55:42 [debug]: len(resps)=1, resp.Reply()=��
2019/11/13 16:55:42 [debug]: call done
[TCP] Int.Add(1, 222) got 223, want: 223
```