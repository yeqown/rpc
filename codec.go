package rpc

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

var (
	_ Codec = &gobCodec{}
)

func init() {
	gob.Register([]*stdRequest{})
	gob.Register([]*stdResponse{})

	gob.Register(&stdRequest{})
	gob.Register(&stdResponse{})
}

// Codec ... to encode and decode
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

// NewGobCodec to generate a new gobCodec instance
func NewGobCodec() Codec {
	// encBuf := bytes.NewBuffer(nil)
	// decBuf := bytes.NewBuffer(nil)

	codec := &gobCodec{
		// encBuf: encBuf,
		// enc:    gob.NewEncoder(encBuf),
		// decBuf: decBuf,
		// dec:    gob.NewDecoder(decBuf),
	}
	return codec
}

// gobCodec using gob to serail
type gobCodec struct {
	// decBuf *bytes.Buffer
	// dec    *gob.Decoder
	// encBuf *bytes.Buffer
	// enc    *gob.Encoder
}

func (g *gobCodec) Encode(argv interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(argv); err != nil {
		return nil, fmt.Errorf("g.enc.Encode(argv) got err: %v", err)
	}

	return buf.Bytes(), nil
}

func (g *gobCodec) Decode(data []byte, out interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	if err := dec.Decode(out); err != nil {
		return fmt.Errorf("[Decode] got err: %v", err)
	}

	return nil
}

// NewRequest ...
func (g *gobCodec) NewResponse(replyv interface{}) Response {
	// resp := new(stdResponse)
	byts, err := g.Encode(replyv)
	if err != nil {
		DebugF("[NewResponse] could not encode replyv=%v, err=%v", replyv, err)
		return nil
	}
	resp := &stdResponse{
		Rply:    byts,
		Err:     "",
		Errcode: Success,
	}

	return resp
}

// ErrResponse .
func (g *gobCodec) ErrResponse(errcode int, err error) Response {
	resp := &stdResponse{
		Err:     err.Error(),
		Errcode: errcode,
	}
	return resp
}

// NewRequest ...
func (g *gobCodec) NewRequest(method string, argv interface{}) Request {
	byts, err := g.Encode(argv)
	if err != nil {
		DebugF("could not encode argv, err=%v", err)
		return nil
	}

	req := &stdRequest{
		Mthd: method,
		Args: byts,
	}

	return req
}

// ReadRequest ...
func (g *gobCodec) ReadRequest(data []byte) ([]Request, error) {
	reqs := make([]Request, 0)
	if err := g.Decode(data, &reqs); err != nil {
		DebugF("[ReadRequest] could not g.Decode(data, stdreqs), err=%v", err)
		return nil, err
	}

	return reqs, nil
	// reqs := make([]Request, len(stdreqs))
	// for idx, v := range stdreqs {
	// 	reqs[idx] = v
	// }

	// return reqs, nil
}

func (g *gobCodec) ReadResponseBody(respBody []byte, rcvr interface{}) error {
	return g.Decode(respBody, rcvr)
}

// ReadRequestBody .
func (g *gobCodec) ReadRequestBody(reqBody []byte, rcvr interface{}) error {
	return g.Decode(reqBody, rcvr)
}

// ReadResponse ...
func (g *gobCodec) ReadResponse(data []byte) ([]Response, error) {
	// DebugF("read response: %v", data)
	resps := make([]Response, 0)
	if err := g.Decode(data, &resps); err != nil {
		return nil, fmt.Errorf("could not decode response: %v", err)
	}

	// resps := make([]Response, len(stdresps))
	// for idx, v := range stdresps {
	// 	resps[idx] = v
	// }
	return resps, nil
}

// EncodeResponses .
func (g *gobCodec) EncodeResponses(v interface{}) ([]byte, error) {
	return g.Encode(v)
}

// EncodeRequests .
func (g *gobCodec) EncodeRequests(v interface{}) ([]byte, error) {
	return g.Encode(v)
}
