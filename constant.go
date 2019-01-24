package rpc

import (
	"fmt"
)

const (
	SUCCESS         = 0
	ParseErr        = -32700 // ParseErr -32700 语法解析错误,服务端接收到无效的json。该错误发送于服务器尝试解析json文本
	InvalidRequest  = -32600 // InvalidRequest -32600 无效请求发送的json不是一个有效的请求对象。
	MethodNotFound  = -32601 // MethodNotFound -32601 找不到方法 该方法不存在或无效
	InvalidParamErr = -32602 // InvalidParamErr -32602 无效的参数 无效的方法参数。
	InternalErr     = -32603 // InternalErr -32603 内部错误 JSON-RPC内部错误。
	// ServerErr       = -32000 // ServerErr -32000 to -32099 Server error服务端错误, 预留用于自定义的服务器错误。
)

type rpcError struct {
	Code   int
	ErrMsg string
}

func (r *rpcError) Error() string {
	return fmt.Sprintf("rpcError(code: %d, errmsg: %s)", r.Code, r.ErrMsg)
}

var errcodeMap = map[int]*rpcError{
	ParseErr:        &rpcError{ParseErr, "ParseErr"},
	InvalidRequest:  &rpcError{InvalidRequest, "InvalidRequest"},
	MethodNotFound:  &rpcError{MethodNotFound, "MethodNotFound"},
	InvalidParamErr: &rpcError{InvalidParamErr, "InvalidParamErr"},
	InternalErr:     &rpcError{InternalErr, "InternalErr"},
	// ServerErr:       "ServerErr",
}
