package json2

import (
	"encoding/json"
	"fmt"

	"github.com/yeqown/rpc"
	"github.com/yeqown/rpc/utils"
)

var (
	_ rpc.Codec = &stdJSONCodec{}
)

// NewStdJSONCodec ...
func NewStdJSONCodec() *stdJSONCodec {
	return &stdJSONCodec{}
}

// stdJSONCodec is suggested for RPC over HTTP
type stdJSONCodec struct{}

func (std *stdJSONCodec) Encode(argv interface{}) ([]byte, error) {
	byts, err := json.Marshal(argv)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal(argv) got err: %v", err)
	}
	return byts, nil
}

func (std *stdJSONCodec) Decode(data []byte, out interface{}) error {
	return json.Unmarshal(data, out)
}

func (std *stdJSONCodec) Response(req rpc.Request, reply []byte, errcode int) rpc.Response {
	resp := &jsonResponse{
		Version: VERSIONCODE,
	}

	if req != nil {
		jsonReq := req.(*jsonRequest)
		resp.ID = jsonReq.ID
	}

	if errcode != rpc.SUCCESS {
		resp.Err = errcodeMap[errcode]
	} else {
		resp.Result = reply
	}

	return resp
}

func (std *stdJSONCodec) Request(method string, argv interface{}) rpc.Request {
	byts, err := std.Encode(argv)
	if err != nil {
		panic(err)
	}

	// if argv is a list
	req := &jsonRequest{
		ID:      utils.RandID(),
		Mthd:    method,
		Args:    byts,
		Version: VERSIONCODE,
	}

	return req
}

func (std *stdJSONCodec) ParseResponse(data []byte) (rpc.Response, error) {
	// TODO: could be jsonResponseArray
	resp := new(jsonResponse)
	if err := std.Decode(data, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (std *stdJSONCodec) ParseRequest(data []byte) (rpc.Request, error) {
	// TODO: could be jsonRequestArray
	req := new(jsonRequest)
	if err := std.Decode(data, req); err != nil {
		return nil, err
	}
	return req, nil
}
