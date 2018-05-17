# rpc (json-rpc 2.0)

remote procedure call over TCP, only support json-rpc 2.0.

### TCP json-rpc Sample

```golang
// test_client.go
package main

import (
	"rpc"
)

type Args struct {
	A int `json:"a"`
	B int `json:"b"`
}

type MultyArgs struct {
	A *Args `json:"aa"`
	B *Args `json:"bb"`
}

type MultyReply struct {
	A int `json:"aa"`
	B int `json:"bb"`
}

func main() {
	c := rpc.NewClient()
	c.DialTCP("127.0.0.1:9999")

	var sum int
	c.Call("1", "Int.Sum", &Args{A: 1, B: 2}, &sum)
	println(sum)

	c.DialTCP("127.0.0.1:9999")
	var reply MultyReply
	c.Call("2", "Int.Multy", &MultyArgs{A: &Args{1, 2}, B: &Args{3, 4}}, &reply)
	println(reply.A, reply.B)
}


```

```golang
// test_server.go
package main

import (
	// "fmt"
	"rpc"
)

type Int int

type Args struct {
	A int `json:"a"`
	B int `json:"b"`
}

func (i *Int) Sum(args *Args, reply *int) error {
	// println("called", args.A, args.B)
	*reply = args.A + args.B
	// *reply = 2
	return nil
}

type MultyArgs struct {
	A *Args `json:"aa"`
	B *Args `json:"bb"`
}

type MultyReply struct {
	A int `json:"aa"`
	B int `json:"bb"`
}

func (i *Int) Multy(args *MultyArgs, reply *MultyReply) error {
	reply.A = (args.A.A * args.A.B)
	reply.B = (args.B.A * args.B.B)
	// fmt.Println(*args.A, *args.B, *reply)
	return nil
}

func main() {
	s := rpc.NewServer()
	mine_int := new(Int)
	s.Register(mine_int)
	s.HandleTCP("127.0.0.1:9999")
}
```

# result

![server](https://raw.githubusercontent.com/yeqown/rpc/master/screenshot/server.png)
![client](https://raw.githubusercontent.com/yeqown/rpc/master/screenshot/client.png)
