package json2

import (
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

type jsonRequest struct {
	ID      string `json:"id"`
	Mthd    string `json:"method"`
	Args    []byte `json:"params"`
	Version string `json:"jsonrpc"`
}

func (j *jsonRequest) Method() string                      { return j.Mthd }
func (j *jsonRequest) Params() []byte                      { return j.Args }
func (j *jsonRequest) CanIter() bool                       { return false }
func (j *jsonRequest) Iter(iterFunc func(req rpc.Request)) { return }

type jsonResponse struct {
	ID      string     `json:"id"`
	Err     *jsonError `json:"error,omitempty"`
	Result  []byte     `json:"result,omitempty"`
	Version string     `json:"jsonrpc"`
}

func (j *jsonResponse) Reply() []byte { return j.Result }
func (j *jsonResponse) Error() error  { return j.Err }
func (j *jsonResponse) ErrCode() int {
	if j.Err == nil {
		return rpc.SUCCESS
	}
	return j.Err.Code
}

// type jsonMultiRequest []jsonRequest

// type jsonMultiResponse []jsonResponse
