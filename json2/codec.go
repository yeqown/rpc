package json2

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/yeqown/rpc"
	"github.com/yeqown/rpc/utils"
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

func (j *jsonCodec) Encode(argv interface{}) ([]byte, error) {
	// byts, err := json.Marshal(argv)
	j.encBuf.Reset()
	err := j.enc.Encode(argv)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal(argv) got err: %v", err)
	}
	byts := j.encBuf.Bytes()
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
		return fmt.Errorf("base64.StdEncoding.DecodeString('%s') got err: %v", string(data), err)
	}
	// println(string(data), len(data), string(base64Dst), len(base64Dst), base64Dst)
	// log.Printf("%v, base64.StdEncoding.DecodedLen(len(data)): %d\n", base64Dst, base64.StdEncoding.DecodedLen(len(data)))
	if _, err := j.decBuf.Write(base64Dst); err != nil {
		return fmt.Errorf("j.decBuf.Write(base64Dst) got err: %v", err)
	}
	return j.dec.Decode(out)
	// return json.Unmarshal(base64Dst, out)
}

func (j *jsonCodec) Response(req rpc.Request, reply []byte, errcode int) rpc.Response {
	resp := &jsonResponse{
		Version: VERSIONCODE,
		ID:      "",
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

func (j *jsonCodec) ParseResponse(data []byte) (resp rpc.Response, err error) {
	resp = new(jsonResponse)
	if err := j.Decode(data, resp); err != nil {

		log.Println("try to decode into jsonResponseArray")
		resp = new(jsonResponseArray)
		if err = j.Decode(data, resp); err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (j *jsonCodec) ParseRequest(data []byte) (req rpc.Request, err error) {
	req = new(jsonRequest)
	if err := j.Decode(data, req); err != nil {

		// try to decode multi
		log.Println("try to decode into jsonRequestArray")
		req = new(jsonRequestArray)
		if err = j.Decode(data, req); err != nil {
			return nil, err
		}
	}
	return req, nil
}

func (j *jsonCodec) MultiSupported() bool { return true }
func (j *jsonCodec) ResponseMulti(resps []rpc.Response) rpc.Response {
	data := make([]*jsonResponse, len(resps))
	for idx := range resps {
		data[idx] = resps[idx].(*jsonResponse)
	}
	return &jsonResponseArray{Data: data}
}
func (j *jsonCodec) RequestMulti(cfgs []*rpc.RequestConfig) rpc.Request {
	reqs := make([]*jsonRequest, len(cfgs))
	for idx := range cfgs {
		byts, _ := j.Encode(cfgs[idx].Args)
		reqs[idx] = &jsonRequest{
			ID:      utils.RandID(),
			Mthd:    cfgs[idx].Method,
			Version: VERSIONCODE,
			Args:    byts,
		}
	}
	return &jsonRequestArray{Data: reqs}
}
