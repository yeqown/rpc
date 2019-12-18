package jsonrpc

import (
	"encoding/json"
	"fmt"

	"github.com/yeqown/rpc"
)

var (
	_ rpc.Request  = &jsonRequest{}
	_ rpc.Response = &jsonResponse{}
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
	ID      string      `json:"id"`
	Mthd    string      `json:"method"`
	Args    interface{} `json:"params"`
	Version string      `json:"jsonrpc"`
}

func (j *jsonRequest) Ident() string  { return j.ID }
func (j *jsonRequest) Method() string { return j.Mthd }
func (j *jsonRequest) Params() []byte {
	byts, err := json.Marshal(j.Args)
	if err != nil {
		panic(err)
	}
	return byts
}

// jsonResponse is jsonCodec response data struct
// and implement the inerface named 'rpc.Response'
type jsonResponse struct {
	ID      string      `json:"id"`
	Err     *jsonError  `json:"error,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Version string      `json:"jsonrpc"`
}

func (j *jsonResponse) SetReqIdent(ident string) { j.ID = ident }
func (j *jsonResponse) Error() error             { return j.Err }
func (j *jsonResponse) Reply() []byte {
	byts, err := json.Marshal(j.Result)
	if err != nil {
		panic(err)
	}
	return byts
}
func (j *jsonResponse) ErrCode() int {
	if j.Err == nil {
		return rpc.Success
	}
	return j.Err.Code
}
