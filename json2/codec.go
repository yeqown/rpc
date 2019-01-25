package json2

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/yeqown/rpc"
	"github.com/yeqown/rpc/utils"
)

var (
	_ rpc.Codec = &jsonCodec{}
)

// NewJSONCodec ...
func NewJSONCodec() *jsonCodec {
	return &jsonCodec{}
}

type jsonCodec struct{}

func (j *jsonCodec) Encode(argv interface{}) ([]byte, error) {
	byts, err := json.Marshal(argv)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal(argv) got err: %v", err)
	}
	base64Dst := make([]byte, base64.StdEncoding.EncodedLen(len(byts)))
	base64.StdEncoding.Encode(base64Dst, byts)
	// println(string(byts), len(byts), string(base64Dst), len(base64Dst))

	return base64Dst, nil
}

func (j *jsonCodec) Decode(data []byte, out interface{}) error {
	// base64Dst := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	var (
		err       error
		base64Dst []byte
	)
	// TOFIX: cannot use base64.StdEncoding.Decode
	if base64Dst, err = base64.StdEncoding.DecodeString(string(data)); err != nil {
		return fmt.Errorf("base64.StdEncoding.Decode(base64Dst, data) got err: %v", err)
	}
	// println(string(data), len(data), string(base64Dst), len(base64Dst), base64Dst)
	// log.Printf("%v, base64.StdEncoding.DecodedLen(len(data)): %d\n", base64Dst, base64.StdEncoding.DecodedLen(len(data)))

	return json.Unmarshal(base64Dst, out)
}

func (j *jsonCodec) Response(req rpc.Request, reply []byte, errcode int) rpc.Response {
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

func (j *jsonCodec) Request(method string, argv interface{}) rpc.Request {
	byts, err := j.Encode(argv)
	if err != nil {
		panic(err)
	}

	req := &jsonRequest{
		ID:      utils.RandID(),
		Mthd:    method,
		Args:    byts,
		Version: VERSIONCODE,
	}

	return req
}

func (j *jsonCodec) ParseResponse(data []byte) (rpc.Response, error) {
	resp := new(jsonResponse)
	if err := j.Decode(data, resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (j *jsonCodec) ParseRequest(data []byte) (rpc.Request, error) {
	req := new(jsonRequest)
	if err := j.Decode(data, req); err != nil {
		return nil, err
	}
	return req, nil
}
