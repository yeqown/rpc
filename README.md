# RPC lib based Golang

In distributed computing, a remote procedure call (RPC) is when a computer program causes a procedure (subroutine) to execute in a different address space (commonly on another computer on a shared network), which is coded as if it were a normal (local) procedure call, without the programmer explicitly coding the details for the remote interaction. That is, the programmer writes essentially the same code whether the subroutine is local to the executing program, or remote. This is a form of client–server interaction (caller is client, executor is server), typically implemented via a request–response message-passing system. In the object-oriented programming paradigm, RPC calls are represented by remote method invocation (RMI). The RPC model implies a level of location transparency, namely that calling procedures is largely the same whether it is local or remote, but usually they are not identical, so local calls can be distinguished from remote calls. Remote calls are usually orders of magnitude slower and less reliable than local calls, so distinguishing them is important.

## Todos

* [x] Codec feature.
* [x] RPC implemention over TCP.
* [ ] JSON RPC(v2) implemention over TCP and HTTP.
* [ ] more test cases.

## Documention

reference to: [godoc](https://godoc.org/github.com/yeqown/rpc)

## Examples

### [RPC example](examples/rpc)
### [JSONRPC example](examples/json2)