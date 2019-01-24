package json2

import (
	"errors"

	"github.com/yeqown/rpc"
)

var (
	_ rpc.Request  = &jsonRequest{}
	_ rpc.Response = &jsonResponse{}
)

type jsonRequest struct {
	ID      string `json:"id"`
	Mthd    string `json:"method"`
	Args    []byte `json:"params"`
	Version string `json:"jsonrpc"`
}

func (j *jsonRequest) Method() string {
	return j.Mthd
}

func (j *jsonRequest) Params() []byte {
	return j.Args
}
func (j *jsonRequest) CanIter() bool {
	return false
}
func (j *jsonRequest) Iter(iterFunc func(req rpc.Request)) {
	return
}

type jsonResponse struct {
	ID      string `json:"id"`
	Err     string `json:"error,omitempty"`
	Result  []byte `json:"result,omitempty"`
	Version string `json:"jsonrpc"`
}

func (j *jsonResponse) Reply() []byte {
	return j.Result
}

func (j *jsonResponse) Error() error {
	if j.Err == "" {
		return nil
	}
	return errors.New(j.Err)
}

type jsonMultiRequest []jsonRequest

type jsonMultiResponse []jsonResponse
