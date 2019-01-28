package json2

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"log"

// 	"github.com/yeqown/rpc"
// 	"github.com/yeqown/rpc/utils"
// )

// var (
// 	_ rpc.Codec = &stdJSONCodec{}
// )

// // NewStdJSONCodec ...
// func NewStdJSONCodec() rpc.Codec {
// 	decbuf := bytes.NewBuffer(nil)
// 	encbuf := bytes.NewBuffer(nil)
// 	c := &stdJSONCodec{
// 		decBuf: decbuf,
// 		dec:    json.NewDecoder(decbuf),
// 		encBuf: encbuf,
// 		enc:    json.NewEncoder(encbuf),
// 	}
// 	c.dec.DisallowUnknownFields()
// 	return c
// }

// // stdJSONCodec is suggested for RPC over HTTP
// type stdJSONCodec struct {
// 	decBuf *bytes.Buffer
// 	dec    *json.Decoder
// 	encBuf *bytes.Buffer
// 	enc    *json.Encoder
// }

// func (std *stdJSONCodec) Encode(argv interface{}) ([]byte, error) {
// 	std.encBuf.Reset()
// 	// std.enc = json.NewEncoder(std.encBuf)
// 	err := std.enc.Encode(argv)
// 	// byts, err := json.Marshal(argv)
// 	if err != nil {
// 		return nil, fmt.Errorf("std.enc.Encode(argv) got err: %v", err)
// 	}
// 	return std.encBuf.Bytes(), nil
// }

// func (std *stdJSONCodec) Decode(data []byte, out interface{}) error {
// 	std.decBuf.Reset()
// 	// std.dec = json.NewDecoder(std.decBuf)
// 	if _, err := std.decBuf.Write(data); err != nil {
// 		return fmt.Errorf("std.decBuf.Write(data) got err: %v", err)
// 	}
// 	return std.dec.Decode(out)
// }

// func (std *stdJSONCodec) Response(req rpc.Request, reply []byte, errcode int) rpc.Response {
// 	resp := &jsonResponse{
// 		Version: VERSIONCODE,
// 	}

// 	if req != nil {
// 		jsonReq := req.(*jsonRequest)
// 		resp.ID = jsonReq.ID
// 	}

// 	if errcode != rpc.SUCCESS {
// 		resp.Err = errcodeMap[errcode]
// 	} else {
// 		resp.Result = reply
// 	}

// 	return resp
// }

// func (std *stdJSONCodec) Request(method string, argv interface{}) rpc.Request {
// 	byts, err := std.Encode(argv)
// 	if err != nil {
// 		panic(err)
// 	}

// 	req := &jsonRequest{
// 		ID:      utils.RandID(),
// 		Mthd:    method,
// 		Args:    byts,
// 		Version: VERSIONCODE,
// 	}

// 	return req
// }

// func (std *stdJSONCodec) ParseResponse(data []byte) (resp rpc.Response, err error) {
// 	resp = new(jsonResponse)
// 	if err = std.Decode(data, resp); err != nil {
// 		// try to decode to jsonResponseArray
// 		resp = new(jsonResponseArray)
// 		if err = std.Decode(data, resp); err != nil {
// 			return nil, err
// 		}
// 	}
// 	return resp, nil
// }

// func (std *stdJSONCodec) ParseRequest(data []byte) (req rpc.Request, err error) {
// 	// log.Printf("(std *stdJSONCodec) ParseRequest(data []byte) input %s", data)
// 	req = new(jsonRequest)
// 	if err = std.Decode(data, req); err != nil {
// 		log.Println("decode single request error, try decode to multi")

// 		// try to decode to jsonRequestArray
// 		req = new(jsonRequestArray)
// 		if err = std.Decode(data, req); err != nil {
// 			return nil, err
// 		}
// 	}
// 	return req, nil
// }

// func (std *stdJSONCodec) MultiSupported() bool {
// 	return true
// }

// func (std *stdJSONCodec) ResponseMulti(resps []rpc.Response) rpc.Response {
// 	data := make([]*jsonResponse, len(resps))
// 	for idx := range resps {
// 		data[idx] = resps[idx].(*jsonResponse)
// 	}
// 	log.Printf("data %v", data)
// 	return &jsonResponseArray{Data: data}
// }

// func (std *stdJSONCodec) RequestMulti(cfgs []*rpc.RequestConfig) rpc.Request {
// 	reqs := make([]*jsonRequest, len(cfgs))
// 	for idx := range cfgs {
// 		byts, _ := std.Encode(cfgs[idx].Args)
// 		reqs[idx] = &jsonRequest{
// 			ID:      utils.RandID(),
// 			Mthd:    cfgs[idx].Method,
// 			Version: VERSIONCODE,
// 			Args:    byts,
// 		}
// 	}
// 	log.Printf("requestMulti[1]: %v, [0]: %v", reqs[1].Args, reqs[0].Args)
// 	return &jsonRequestArray{Data: reqs}
// }
