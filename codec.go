package rpc

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"net"
)

var (
	_ Codec = &gobCodec{}
)

// Codec ... to encode and decode
// for client to encode request and decode response
// for server to encode response den decode request
type Codec interface {
	Encode(argv interface{}) ([]byte, error)
	Decode(data []byte, argv interface{}) error

	Response(conn net.Conn, req Request, reply []byte, err error) error
	ParseResponse(respBody []byte) (Response, error)

	Request(conn net.Conn, method string, argv interface{}) ([]byte, error)
	ParseRequest(data []byte) (Request, error)
}

func newGobCodec() *gobCodec {
	encBuf := bytes.NewBuffer(nil)
	decBuf := bytes.NewBuffer(nil)

	codec := &gobCodec{
		encBuf: encBuf,
		enc:    gob.NewEncoder(encBuf),
		decBuf: decBuf,
		dec:    gob.NewDecoder(decBuf),
	}
	return codec
}

// gobCodec using gob to serail
type gobCodec struct {
	decBuf *bytes.Buffer
	dec    *gob.Decoder
	encBuf *bytes.Buffer
	enc    *gob.Encoder
}

func (g *gobCodec) Encode(argv interface{}) ([]byte, error) {
	g.encBuf.Reset()
	g.enc = gob.NewEncoder(g.encBuf)

	if err := g.enc.Encode(argv); err != nil {
		return nil, fmt.Errorf("g.enc.Encode(argv) got err: %v", err)
	}

	src := g.encBuf.Bytes()
	bas64Dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(bas64Dst, src)

	return bas64Dst, nil
}

func (g *gobCodec) Decode(data []byte, out interface{}) error {
	g.decBuf.Reset()
	g.dec = gob.NewDecoder(g.decBuf)

	base64Dst := make([]byte, base64.StdEncoding.DecodedLen(len(data)))
	if _, err := base64.StdEncoding.Decode(base64Dst, data); err != nil {
		return fmt.Errorf("hex.Decode(base64Dst, data) got err: %v", err)
	}

	if _, err := g.decBuf.Write(base64Dst); err != nil {
		return fmt.Errorf("g.decBuf.Write(base64Dst) got err: %v", err)
	}

	if err := g.dec.Decode(out); err != nil {
		return fmt.Errorf("g.dec.Decode(out) got err: %v", err)
	}
	return nil
}

// Response ...
func (g *gobCodec) Response(conn net.Conn, req Request, reply []byte, err error) error {
	resp := new(defaultResponse)
	if err != nil {
		resp.Err = err.Error()
	} else {
		resp.Rply = reply
	}

	return WriteServerTCP(conn, g, resp)
}

// Request ...
func (g *gobCodec) Request(conn net.Conn, method string, argv interface{}) ([]byte, error) {
	byts, err := g.Encode(argv)
	if err != nil {
		return nil, err
	}

	req := &defaultRequest{
		Mthd: method,
		Args: byts,
	}

	return WriteClientTCP(conn, g, req)
}

// ParseRequest ...
func (g *gobCodec) ParseRequest(data []byte) (Request, error) {
	req := new(defaultRequest)
	if err := g.Decode(data, req); err != nil {
		return nil, err
	}
	return req, nil
}

// ParseResponse ...
func (g *gobCodec) ParseResponse(data []byte) (Response, error) {
	resp := new(defaultResponse)
	if err := g.Decode(data, resp); err != nil {
		return nil, fmt.Errorf("could not decode response: %v", err)
	}
	return resp, nil
}
