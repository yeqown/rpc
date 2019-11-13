package jsonrpc

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/yeqown/rpc"
)

var (
	_ rpc.Codec = &jsonCodec{}
)

// NewJSONCodec ...
func NewJSONCodec() rpc.Codec {
	decbuf := bytes.NewBuffer(nil)
	encbuf := bytes.NewBuffer(nil)
	c := &jsonCodec{
		decBuf: decbuf,
		dec:    json.NewDecoder(decbuf),
		encBuf: encbuf,
		enc:    json.NewEncoder(encbuf),
	}
	c.dec.DisallowUnknownFields()
	return c
}

type jsonCodec struct {
	decBuf *bytes.Buffer
	dec    *json.Decoder
	encBuf *bytes.Buffer
	enc    *json.Encoder
}

func (j *jsonCodec) encode(argv interface{}) ([]byte, error) {
	j.encBuf.Reset()
	err := j.enc.Encode(argv)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal(argv) got err: %v", err)
	}
	// log.Printf("j.encBuf cap: %d", j.encBuf.Cap())
	return j.encBuf.Bytes(), nil
}

func (j *jsonCodec) decode(data []byte, out interface{}) error {
	j.decBuf.Truncate(0)
	var (
		err error
	)
	if _, err = j.decBuf.Write(data); err != nil {
		return fmt.Errorf("j.decBuf.Write(data) got err: %v", err)
	}
	// log.Printf("j.decBuf cap: %d", j.decBuf.Cap())
	return j.dec.Decode(out)
}

func (j *jsonCodec) NewResponse(reply interface{}) rpc.Response {
	resp := &jsonResponse{
		Version: VERSIONCODE, // this is a dead string '2.0'
		ID:      "",          // this will be set in server side
		Result:  reply,       // interface{}
	}

	return resp
}

func (j *jsonCodec) ErrResponse(errcode int, err error) rpc.Response {
	errmsg := ""
	if err != nil {
		errmsg = err.Error()
	}

	return &jsonResponse{
		Err: &jsonError{
			Code:    errcode,
			Message: errmsg,
		},
		Version: VERSIONCODE,
	}
}

func (j *jsonCodec) NewRequest(method string, argv interface{}) rpc.Request {
	req := &jsonRequest{
		ID:      randid(),
		Mthd:    method,
		Args:    argv,
		Version: VERSIONCODE,
	}

	return req
}

func (j *jsonCodec) ReadResponse(data []byte) (resps []rpc.Response, err error) {
	jsonResps := make([]*jsonResponse, 0)
	if err := j.decode(data, &jsonResps); err != nil {
		// try to decode multi
		rpc.DebugF("try to decode into jsonResponseArray, err=%v", err)
		resp := new(jsonResponse)
		if err := j.decode(data, resp); err != nil {
			rpc.DebugF("try to decode into jsonResponse failed, err=%v", err)
			return nil, err
		}
		resps = append(resps, resp)
		return resps, nil
	}

	for _, jsonResp := range jsonResps {
		resps = append(resps, jsonResp)
	}

	return resps, nil
}

func (j *jsonCodec) ReadRequest(data []byte) (reqs []rpc.Request, err error) {
	jsonReqs := make([]*jsonRequest, 0)
	if err := j.decode(data, &jsonReqs); err != nil {
		// try to decode multi
		rpc.DebugF("try to decode into jsonRequestArray, err=%v", err)
		req := new(jsonRequest)
		if err := j.decode(data, req); err != nil {
			rpc.DebugF("try to decode into jsonRequest failed, err=%v", err)
			return nil, err
		}
		reqs = append(reqs, req)
		return reqs, nil
	}

	for _, jsonReq := range jsonReqs {
		reqs = append(reqs, jsonReq)
	}

	return reqs, nil
}

func (j *jsonCodec) ReadRequestBody(data []byte, rcvr interface{}) error {
	return j.decode(data, rcvr)
}

func (j *jsonCodec) ReadResponseBody(data []byte, rcvr interface{}) error {
	return j.decode(data, rcvr)
}

func (j *jsonCodec) EncodeRequests(v interface{}) ([]byte, error) {
	return j.encode(v)
}

func (j *jsonCodec) EncodeResponses(v interface{}) ([]byte, error) {
	return j.encode(v)
}
