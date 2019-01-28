# RPC lib based Golang
[![Go Report Card](https://goreportcard.com/badge/github.com/yeqown/rpc)](https://goreportcard.com/report/github.com/yeqown/rpc) [![GoReportCard](https://godoc.org/github.com/yeqown/rpc?status.svg)](https://godoc.org/github.com/yeqown/rpc)

In distributed computing, a remote procedure call (RPC) is when a computer program causes a procedure (subroutine) to execute in a different address space (commonly on another computer on a shared network), which is coded as if it were a normal (local) procedure call, without the programmer explicitly coding the details for the remote interaction. That is, the programmer writes essentially the same code whether the subroutine is local to the executing program, or remote. This is a form of client–server interaction (caller is client, executor is server), typically implemented via a request–response message-passing system. In the object-oriented programming paradigm, RPC calls are represented by remote method invocation (RMI). The RPC model implies a level of location transparency, namely that calling procedures is largely the same whether it is local or remote, but usually they are not identical, so local calls can be distinguished from remote calls. Remote calls are usually orders of magnitude slower and less reliable than local calls, so distinguishing them is important.

## Todos

* [x] Codec feature.
* [x] RPC implemention over TCP.
* [x] RPC implemention over HTTP.
* [ ] JSON RPC(v2) implemention over TCP and HTTP.
* [ ] more test cases.
* [ ] compatible with JSON RPC 1.0

## Documention

### 1. all API reference to [godoc](https://godoc.org/github.com/yeqown/rpc)

### 2. usage over TCP:

#### 2.1 server side

```go
func main() {
	srv := rpc.NewServerWithCodec(nil)
	// srv.Register(new(Int)) will register all exported methods of type `Int`
	srv.RegisterName(new(Int), "Add")
	// srv.Start(tcpAddr, httpAddr) will serve both TCP and HTTP
	srv.ServeTCP(tcpAddr)
}
```

#### 2.2 client side

```go

func main() {
	c := rpc.NewClientWithCodec(nil, "127.0.0.1:9998", "")
	var (
		sum  int
		args = &Args{A: 1, B: 222}
	)
	if err := c.CallOverTCP("Int.Add", args, &sum); err != nil {
		println("got err: ", err.Error())
	}

	fmt.Printf("[TCP] Int.Add(%d, %d) got %d, want: %d\n", args.A, args.B, sum, args.A+args.B)
}
```

### 3. usage over HTTP:

#### 3.1 server side

```go
func main() {
	srv := rpc.NewServerWithCodec(nil)
	// srv.Register(new(Int)) will register all exported methods of type `Int`
	srv.RegisterName(new(Int), "Add")
	srv.ListenAndServe(httpAddr)
	// srv.Start(tcpAddr, httpAddr) will serve both TCP and HTTP
}
```

#### 3.2 client side

```go
func main() {
	// prototype rpc.NewClientWithCodec(codec Codec, tcpAddr string, httpAddr string)
	// if codec is nil will use default gobCodec, tcpAddr or httpAddr is empty only when
	// you are sure about it will never be used, otherwise it panic while using some functions.
	c := rpc.NewClientWithCodec(nil, "", "127.0.0.1:9999")
	var (
		sum  int
		args = &Args{A: 1111, B: 222}
	)
	if err := c.CallOverHTTP("Int.Add", args, &sum); err != nil {
		println("got err: ", err.Error())
	}
	fmt.Printf("[HTTP] Int.Add(%d, %d) got %d, want: %d\n", args.A, args.B, sum, args.A+args.B)
}
```

## Examples

### [RPC example](examples/rpc)
### [JSON RPC 2.0 example](examples/json2)