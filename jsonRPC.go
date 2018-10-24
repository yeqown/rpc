package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

const (
	// JSONRPCVER const version code of JSONRPC
	JSONRPCVER = "2.0"

	// MaxMultiRequest count
	MaxMultiRequest = 10

	// ParseErr -32700 语法解析错误,服务端接收到无效的json。该错误发送于服务器尝试解析json文本
	ParseErr = -32700

	// InvalidRequest -32600 无效请求发送的json不是一个有效的请求对象。
	InvalidRequest = -32600

	// MethodNotFound -32601 找不到方法 该方法不存在或无效
	MethodNotFound = -32601

	// InvalidParamErr -32602 无效的参数 无效的方法参数。
	InvalidParamErr = -32602

	// InternalErr -32603 内部错误 JSON-RPC内部错误。
	InternalErr = -32603

	// ServerErr -32000 to -32099 Server error服务端错误, 预留用于自定义的服务器错误。
	ServerErr = -32000
)

var _messages = map[int]string{
	ParseErr:        "ParseErr",
	InvalidRequest:  "InvalidRequest",
	MethodNotFound:  "MethodNotFound",
	InvalidParamErr: "InvalidParamErr",
	InternalErr:     "InternalErr",
	ServerErr:       "ServerErr",
}

// Response is server send a response to client,
// and client parse to this
type Response struct {
	ID      string      `json:"id"`
	Error   *JsonrpcErr `json:"error"`
	Result  interface{} `json:"result"`
	Jsonrpc string      `json:"jsonrpc"`
}

// NewResponse ...
func NewResponse(id string, result interface{}, err *JsonrpcErr) *Response {
	if err != nil {
		id = ""
	}
	return &Response{
		ID:      id,
		Error:   err,
		Result:  result,
		Jsonrpc: JSONRPCVER,
	}
}

// JsonrpcErr while dealing with rpc request got err,
// must return this.
type JsonrpcErr struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (j *JsonrpcErr) Error() string {
	bs, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

// NewJsonrpcErr ...
func NewJsonrpcErr(code int, message string, data interface{}) *JsonrpcErr {
	if message == "" {
		message = _messages[code]
	}
	return &JsonrpcErr{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

// Request Client send request to server,
// and server also parse request into this
type Request struct {
	ID      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Jsonrpc string      `json:"jsonrpc"`
}

// NewRequest ...
func NewRequest(id string, params interface{}, method string) *Request {
	return &Request{
		ID:      id,
		Params:  params,
		Method:  method,
		Jsonrpc: JSONRPCVER,
	}
}

// encodeRequest
func encodeRequest(req *Request) []byte {
	bs, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	return bs
}

// encodeMultiRequest
func encodeMultiRequest(reqs []*Request) []byte {
	bs, err := json.Marshal(reqs)
	if err != nil {
		panic(err)
	}
	return bs
}

// parseRequest parse request data []byte into []*Request
// try to parse multi request first,
// if the func gets any err then parse the body to single request
func parseRequest(bs []byte) ([]*Request, error) {
	mr := make([]*Request, 0, MaxMultiRequest)
	// println(string(bs))
	if err := json.Unmarshal(bs, &mr); err != nil {
		log.Println("parseMultiRequest err:", err.Error())
		goto ParseSingleReq
	}
	return mr, nil

ParseSingleReq:
	r := new(Request)
	if err := json.Unmarshal(bs, r); err != nil {
		errmsg := "ParseSingleReq err: " + err.Error()
		log.Println(errmsg)
		return mr, errors.New(errmsg)
	}
	mr = append(mr, r)
	return mr, nil
}

// encodeResponse encode *Response into []byte
// to send to client
func encodeResponse(resp *Response) []byte {
	bs, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	return bs
}

// encodeMultiResponse encode []*Repsonse into []byte
// then can be send mulit response to client
func encodeMultiResponse(resps []*Response) []byte {
	bs, err := json.Marshal(resps)
	if err != nil {
		panic(err)
	}
	return bs
}

// parseResponse ... parse []byte into *Response
func parseResponse(s string) *Response {
	resp := new(Response)
	if err := json.Unmarshal([]byte(s), resp); err != nil {
		errmsg := fmt.Sprintf("%v recvived!", s)
		log.Println(errmsg)
		panic(err)
	}
	return resp
}

// parseMultiResponse ...
func parseMultiResponse(s string) []*Response {
	resps := make([]*Response, 0)
	if err := json.Unmarshal([]byte(s), &resps); err != nil {
		errmsg := fmt.Sprintf("recvived! and parese err: %v", err)
		log.Println(errmsg)
		panic(errmsg)
	}
	return resps
}
