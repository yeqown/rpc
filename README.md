# RPC (JSON-RPC 2.0)

remote procedure call over TCP and HTTP(Post & json-body only), only support JSON-RPC 2.0.

> Note: Currently, tests are inadequate 0.0, I'll do this quickly. And this lib is too slow ~

### Doc

ref to:
[https://godoc.org/github.com/yeqown/rpc](https://godoc.org/github.com/yeqown/rpc)

### Sample Code

link to [code](sample/server/client.go)
```golang
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

// MultyArgs ...
type MultyArgs struct {
	A *Args `json:"aa"`
	B *Args `json:"bb"`
}

// MultyReply ...
type MultyReply struct {
	A int `json:"aa"`
	B int `json:"bb"`
}

func main() {
	c := rpc.NewClient()
	c.DialTCP("127.0.0.1:9999")
	testAdd(c)

	// after c.Call the conn will be closed
	c.DialTCP("127.0.0.1:9999")
	testMultyReply(c)

	c.DialTCP("127.0.0.1:9999")
	testMultiParamAdd(c)
}

func testAdd(c *rpc.Client) {
	var (
		sum     = 0
		wantSum = 3
	)
	c.Call("1", "Int.Sum", &Args{A: 1, B: 2}, &sum)
	if !reflect.DeepEqual(sum, wantSum) {
		err := fmt.Errorf("Int.Sum Result %d not equal to %d", sum, wantSum)
		fmt.Println(err)
		return
	}
	fmt.Println("testAdd passed")
}

func testMultiParamAdd(c *rpc.Client) {
	var (
		params  = make([]*Args, 3)
		sum     = make([]*int, 3)
		wantSum = make([]*int, 3)
	)

	for i := 0; i < 3; i++ {
		params[i] = &Args{A: i, B: i * 2}
		wantSum[i] = new(int)
		sum[i] = new(int)
		*(wantSum[i]) = (i + i*2)
		// allocate the mem for reply , or cannot set the Response.Result to reply
	}

	c.CallMulti("Int.Sum", &params, &sum)
	// fmt.Printf("%v, %v", sum, wantSum)
	// for _, v := range sum {
	// 	fmt.Println(*v)
	// }
	if !reflect.DeepEqual(sum, wantSum) {
		err := fmt.Errorf("Int.Sum Result %v not equal to %v", sum, wantSum)
		fmt.Println(err)
		return
	}
	fmt.Println("testMultiParamAdd passed")
}

func testMultyReply(c *rpc.Client) {
	var (
		reply     MultyReply
		wantReply = MultyReply{
			A: 2,
			B: 12,
		}
	)
	c.Call("2", "Int.Multy", &MultyArgs{A: &Args{1, 2}, B: &Args{3, 4}}, &reply)
	if !reflect.DeepEqual(reply, wantReply) {
		err := fmt.Errorf("Int.Multy Result %v not equal to %v", reply, wantReply)
		fmt.Println(err)
		return
	}
	fmt.Println("testMultyReply passed")
}

```

link to [code](sample/server/server.go)
```golang
package main

import (
	"net/http"

	"github.com/yeqown/rpc"
)

// Int ... custom type for JSON-RPC test
type Int int

// Args ... for Sum Method
type Args struct {
	A int `json:"a"`
	B int `json:"b"`
}

// Sum ...
func (i *Int) Sum(args *Args, reply *int) error {
	// println("called", args.A, args.B)
	*reply = args.A + args.B
	// *reply = 2
	return nil
}

// MultyArgs ... from Multy Int.Method
type MultyArgs struct {
	A *Args `json:"aa"`
	B *Args `json:"bb"`
}

// MultyReply ...
type MultyReply struct {
	A int `json:"aa"`
	B int `json:"bb"`
}

// Multy ... times params
func (i *Int) Multy(args *MultyArgs, reply *MultyReply) error {
	reply.A = (args.A.A * args.A.B)
	reply.B = (args.B.A * args.B.B)
	// fmt.Println(*args.A, *args.B, *reply)
	return nil
}

func main() {
	s := rpc.NewServer()
	i := new(Int)
	s.Register(i)
	go s.HandleTCP("127.0.0.1:9999")

	// to support http Request
	http.ListenAndServe(":9998", s)
}

```

# Sample Running Shotting

![server](https://raw.githubusercontent.com/yeqown/rpc/master/screenshot/server.png)
![client](https://raw.githubusercontent.com/yeqown/rpc/master/screenshot/client.png)
![http-support](https://raw.githubusercontent.com/yeqown/rpc/master/screenshot/http-support.png)
