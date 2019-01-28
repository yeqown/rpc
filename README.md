# RPC lib based Golang
[![Go Report Card](https://goreportcard.com/badge/github.com/yeqown/rpc)](https://goreportcard.com/report/github.com/yeqown/rpc) [![GoReportCard](https://godoc.org/github.com/yeqown/rpc?status.svg)](https://godoc.org/github.com/yeqown/rpc)

In distributed computing, a remote procedure call (RPC) is when a computer program causes a procedure (subroutine) to execute in a different address space (commonly on another computer on a shared network), which is coded as if it were a normal (local) procedure call, without the programmer explicitly coding the details for the remote interaction. That is, the programmer writes essentially the same code whether the subroutine is local to the executing program, or remote. This is a form of client–server interaction (caller is client, executor is server), typically implemented via a request–response message-passing system. In the object-oriented programming paradigm, RPC calls are represented by remote method invocation (RMI). The RPC model implies a level of location transparency, namely that calling procedures is largely the same whether it is local or remote, but usually they are not identical, so local calls can be distinguished from remote calls. Remote calls are usually orders of magnitude slower and less reliable than local calls, so distinguishing them is important.

## Todos

* [x] Codec feature.
* [x] RPC implemention over TCP.
* [x] RPC implemention over HTTP.
* [x] JSON RPC(v2) implemention over TCP and HTTP.
* [ ] more test cases.
* [ ] compatible with JSON RPC 1.0

## Documention

#### About interface `Codec`

`Codec` is a interface contains fields:
```go
// Codec to encode and decode
// for client to encode request and decode response
// for server to encode response den decode request
type Codec interface {
	// Encode an interface value into []byte
	Encode(argv interface{}) ([]byte, error)

	// Decode encoded data([]byte) back to an interface which the origin data belongs to
	Decode(data []byte, argv interface{}) error

	// generate a single Response with needed params
	Response(req Request, reply []byte, errcode int) Response

	// parse encoded data into a Response
	ParseResponse(respBody []byte) (Response, error)

	// generate a single Request with needed params
	Request(method string, argv interface{}) Request

	// parse encoded data into a Request
	ParseRequest(data []byte) (Request, error)

	// if MultiSupported return true means, can provide funcs
	// ResponseMulti, ParseResponseMulti, RequestMulti, ParseRequestMulti
	MultiSupported() bool

	// generate a Response which cann support Iter(iterator interface)
	ResponseMulti(resps []Response) Response

	// generate a Request which cann support Iter(iterator interface)
	RequestMulti(cfgs []*RequestConfig) Request
}
```

#### About struct `RequestConfig`

`RequestConfig` is a data structure to call multiple request config, it contains:
```go
// RequestConfig ... to support request multi
type RequestConfig struct {
	// Method that called by client, if not existed will recv an err.
	Method string

	// Args should be params pointer type
	Args interface{}

	// Reply shoule be result pointer type
	Reply interface{}
}
```

use this config while needing send a multi request at one time, just like this:

```go
cfgs := []*rpc.RequestConfig{
	&rpc.RequestConfig{
		Method: "Int.Add",
		Args:   &Args{10, 1909},
		Reply:  &Result{},
	},
	&rpc.RequestConfig{
		Method: "Int.Sum",
		Args:   &Args{21312, 1909},
		Reply:  &Result{},
	},
}
// c.CallOverTCPMulti(cfgs) is ok too
if err := c.CallOverHTTPMulti(cfgs); err != nil {
	log.Printf("c.CallOverHTTPMulti client got err: %v", err)
}
```

#### Server Side API


`Register` will register all exported method of `rcvr`
```go
func (s *Server) Register(rcvr interface{})
```

`RegisterName` only register the appointed `method` of `rcvr`
```go
func (s *Server) RegisterName(rcvr interface{}, methodName string)
```

`ServeTCP` to run *TCP* server
```go
func (s *Server) ServeTCP(tcpAddr string)
```

`ServeHTTP` to implement `http.Handler` interface so that the server can serve with *HTTP* request
```go
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request)
```

`ListenAndServe` to run an http server and handle each request
```go
func (s *Server) ListenAndServe(httpAddr string)
```

`Start` to run TCP and HTTP server, if any `addr` is empty, the server whose addr is empty will not run
```go
func (s *Server) Start(tcpAddr, httpAddr string)
```

#### Client Side API

`CallOverTCP` call `method` with `argv` and return value into `reply` over *TCP*
```go
func (c *Client) CallOverTCP(method string, argv, reply interface{})
```

`CallOverTCPMulti` send multi request to server configed by `cfgs`, get result form `cfg.Reply`
```go
func (c *Client) CallOverTCPMulti(cfgs []*RequestConfig)
```

`CallOverHTTP` work like `CallOverTCP`, the difference is over *HTTP* rather than *TCP*
```go
func (c *Client) CallOverHTTP(method string, argv, reply interface{})
```

`CallOverHTTPMulti` works like `CallOverTCPMulti`
```go
func (c *Client) CallOverHTTPMulti(tcpAddr, httpAddr string)
```

`Close` `c.conn` (tcp connection)
```go
func (c *Client) Close()
```

## Examples

### [RPC example](examples/rpc)
### [JSON RPC example](examples/json2)