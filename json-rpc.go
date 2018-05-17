package rpc

import (
	"encoding/json"
)

const (
	JSON_RPC_VER    = "2.0"
	MaxMultiRequest = 10

	ParseErr        = -32700 // -32700 语法解析错误,服务端接收到无效的json。该错误发送于服务器尝试解析json文本
	InvalidRequest  = -32600 // -32600 无效请求发送的json不是一个有效的请求对象。
	MethodNotFound  = -32601 // -32601 找不到方法	该方法不存在或无效
	InvalidParamErr = -32602 // -32602 无效的参数	无效的方法参数。
	InternalErr     = -32603 // -32603 内部错误 JSON-RPC内部错误。
	ServerErr       = -32000 // -32000 to -32099	Server error服务端错误	预留用于自定义的服务器错误。
)

var _messages = map[int]string{
	ParseErr:        "ParseErr",
	InvalidRequest:  "InvalidRequest",
	MethodNotFound:  "MethodNotFound",
	InvalidParamErr: "InvalidParamErr",
	InternalErr:     "InternalErr",
	ServerErr:       "ServerErr",
}

// server send a response to client,
// and client parse to this
type Response struct {
	ID      string      `json:"id"`
	Error   *JsonrpcErr `json:"error"`
	Result  interface{} `json:"result"`
	Jsonrpc string      `json:"jsonrpc"`
}

func NewResponse(id string, result interface{}, err *JsonrpcErr) *Response {
	if err != nil {
		id = ""
	}
	return &Response{
		ID:      id,
		Error:   err,
		Result:  result,
		Jsonrpc: JSON_RPC_VER,
	}
}

// JsonrpcErr while dealing with rpc request got err,
// must return this.
type JsonrpcErr struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (je *JsonrpcErr) Error() string {
	bs, err := json.Marshal(je)
	if err != nil {
		panic(err)
	}
	return string(bs)
}

// NewJsonrpcErr
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

// MultiRequest
type MultiRequest []*Request

// Client send request to server,
// and server also parse request into this
type Request struct {
	ID      string      `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	Jsonrpc string      `json:"jsonrpc"`
}

func NewRequest(id string, params interface{}, method string) *Request {
	return &Request{
		ID:      id,
		Params:  params,
		Method:  method,
		Jsonrpc: JSON_RPC_VER,
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
func encodeMultiRequest(reqs *MultiRequest) []byte {
	bs, err := json.Marshal(reqs)
	if err != nil {
		panic(err)
	}
	return bs
}

// parseRequest parse request data string into request
func parseRequest(s string) MultiRequest {
	mr := make(MultiRequest, 0, MaxMultiRequest)
	if err := json.Unmarshal([]byte(s), &mr); err != nil {
		// println("ParseMultiReq err:", err.Error())
		goto ParseSingleReq
	}
	return mr

ParseSingleReq:
	r := new(Request)
	if err := json.Unmarshal([]byte(s), r); err != nil {
		println("ParseSingleReq err: ", err.Error())
		return mr
	}
	mr = append(mr, r)
	return mr
}

// encodeRepsonse
func encodeResponse(resp *Response) []byte {
	bs, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	return bs
}

// encodeMultiResponse
func encodeMultiResponse(resps []*Response) []byte {
	bs, err := json.Marshal(resps)
	if err != nil {
		panic(err)
	}
	return bs
}

// parseResponse
func parseResponse(s string) *Response {
	reps := new(Response)
	if err := json.Unmarshal([]byte(s), reps); err != nil {
		panic(err)
	}
	return reps
}
