package json2

import (
	"fmt"

	"github.com/yeqown/rpc"
)

var (
	_ rpc.Request  = &jsonRequest{}
	_ rpc.Response = &jsonResponse{}
	_ rpc.Request  = &jsonRequestArray{}
	_ rpc.Response = &jsonResponseArray{}
)

type jsonError struct {
	Code    int         `json:"errcode"`
	Message string      `json:"errmsg"`
	Data    interface{} `json:"data,omitempty"`
}

func (j *jsonError) Error() string {
	return fmt.Sprintf("jsonError(code: %d, message: %s)", j.Code, j.Message)
}

// jsonRequest is jsonCodec response data struct
// and implement the inerface named 'rpc.Request'
type jsonRequest struct {
	ID      string `json:"id"`
	Mthd    string `json:"method"`
	Args    []byte `json:"params"`
	Version string `json:"jsonrpc"`
}

func (j *jsonRequest) GetID() string                 { return j.ID }
func (j *jsonRequest) Method() string                { return j.Mthd }
func (j *jsonRequest) Params(codec rpc.Codec) []byte { return j.Args }
func (j *jsonRequest) HasNext() bool                 { return false }
func (j *jsonRequest) Next() interface{}             { return nil }

// jsonResponse is jsonCodec response data struct
// and implement the inerface named 'rpc.Response'
type jsonResponse struct {
	ID      string     `json:"id"`
	Err     *jsonError `json:"error,omitempty"`
	Result  []byte     `json:"result,omitempty"`
	Version string     `json:"jsonrpc"`
}

func (j *jsonResponse) GetID() string                { return j.ID }
func (j *jsonResponse) Reply(codec rpc.Codec) []byte { return j.Result }
func (j *jsonResponse) Error() error {
	if j.Err == nil {
		return nil
	}
	return errcodeMap[j.Err.Code]
}
func (j *jsonResponse) ErrCode() int {
	if j.Err == nil {
		return rpc.SUCCESS
	}
	return j.Err.Code
}
func (j *jsonResponse) HasNext() bool     { return false }
func (j *jsonResponse) Next() interface{} { return nil }

// jsonRequestArray implements Request interface
type jsonRequestArray struct {
	cur  int
	Data []*jsonRequest
}

func (reqa *jsonRequestArray) Method() string { return "" }
func (reqa *jsonRequestArray) Params(codec rpc.Codec) []byte {
	byts, err := codec.Encode(reqa.Data)
	if err != nil {
		panic(err)
	}
	return byts
}
func (reqa *jsonRequestArray) HasNext() bool {
	return reqa.cur < len(reqa.Data)
}
func (reqa *jsonRequestArray) Next() interface{} {
	v := reqa.Data[reqa.cur]
	reqa.cur++
	return v
}

// jsonResponseArray implements Response interface
type jsonResponseArray struct {
	cur  int
	Data []*jsonResponse
}

func (respa *jsonResponseArray) Reply(codec rpc.Codec) []byte {
	byts, err := codec.Encode(respa.Data)
	if err != nil {
		panic(err)
	}
	return byts
}

func (respa *jsonResponseArray) Error() error {
	return nil
}

func (respa *jsonResponseArray) ErrCode() int {
	return rpc.SUCCESS
}

func (respa *jsonResponseArray) HasNext() bool {
	return respa.cur < len(respa.Data)
}

func (respa *jsonResponseArray) Next() interface{} {
	v := respa.Data[respa.cur]
	respa.cur++
	return v
}
