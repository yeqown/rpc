# RPC lib based Golang
[![Go Report Card](https://goreportcard.com/badge/github.com/yeqown/rpc)](https://goreportcard.com/report/github.com/yeqown/rpc) [![](https://godoc.org/github.com/yeqown/rpc?status.svg)](https://godoc.org/github.com/yeqown/rpc)

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
	ServerCodec
	ClientCodec
}


// ServerCodec .
// parse request and write response to client.
type ServerCodec interface {
	// parse encoded data into a Request
	ReadRequest(data []byte) ([]Request, error)
	// ReadRequestBody parse params
	ReadRequestBody(reqBody []byte, rcvr interface{}) error
	// generate a single Response with needed params
	NewResponse(replyv interface{}) Response
	// ErrResponse to generate a Reponse contains error
	ErrResponse(errcode int, err error) Response
	// EncodeResponses .
	EncodeResponses(v interface{}) ([]byte, error)
}

// ClientCodec .
// prase response and write request to server
type ClientCodec interface {
	// generate a single NewRequest with needed params
	NewRequest(method string, argv interface{}) Request
	// EncodeRequests .
	EncodeRequests(v interface{}) ([]byte, error)
	// parse encoded data into a Response
	ReadResponse(data []byte) ([]Response, error)
	// ReadResponseBody .
	ReadResponseBody(respBody []byte, rcvr interface{}) error
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
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.NewRequest)
```

`ListenAndServe` to run an http server and handle each request
```go
func (s *Server) ListenAndServe(httpAddr string)
```

#### Client Side API

`Call` to send a `RPC` request to server and recv response
```go
func (c *Client) Call(method string, args, reply interface{}) error
```

`Close` `c.conn` (tcp connection)
```go
func (c *Client) Close()
```

## Examples

### [RPC gob](examples/rpc)
### [JSON RPC](examples/json2)
### [JSON RPC Array](examples/json2-array)