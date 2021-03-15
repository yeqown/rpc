package jsonrpc

import (
	"testing"

	"github.com/yeqown/rpc"
)

func Test_jsonCodecRequest(t *testing.T) {
	type Args struct {
		A int `json:"a"`
		B int `json:"b"`
	}

	codec := NewJSONCodec().(*jsonCodec)
	req := codec.NewRequest("Int.Sum", &Args{10, 11})
	dat, err := codec.EncodeRequests([]rpc.Request{req})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	reqs, err := codec.ReadRequest(dat)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(reqs) != 1 {
		t.Error("invalid reqs length")
		t.FailNow()
	}

	argsRcvr := new(Args)
	if err := codec.ReadRequestBody(reqs[0].Params(), argsRcvr); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if argsRcvr.A != 10 || argsRcvr.B != 11 {
		t.Error("not equal args")
		t.FailNow()
	}
}

func Test_jsonCodecResponse(t *testing.T) {
	type Args struct {
		A int `json:"a"`
		B int `json:"b"`
	}

	codec := NewJSONCodec().(*jsonCodec)
	resp := codec.NewResponse(&Args{10, 11})

	dat, err := codec.EncodeResponses([]rpc.Response{resp})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	resps, err := codec.ReadResponse(dat)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(resps) != 1 {
		t.Error("invalid resps length")
		t.FailNow()
	}

	argsRcvr := new(Args)
	if err := codec.ReadRequestBody(resps[0].Reply(), argsRcvr); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if argsRcvr.A != 10 || argsRcvr.B != 11 {
		t.Error("not equal args")
		t.FailNow()
	}
}

// https://github.com/yeqown/rpc/issues/5
func Test_jsonCodec_Issues(t *testing.T) {
	codec := NewJSONCodec().(*jsonCodec)

	m := map[string]interface{}{}
	dataWrong := `{"a": 1,}`
	err1 := codec.decode([]byte(dataWrong), &m)
	t.Logf("err1: %v", err1)
	if err1 == nil {
		t.Fatalf("err1 should not be nil")
	}

	dataCorrect := `{"a": 1}`
	err2 := codec.decode([]byte(dataCorrect), &m)
	t.Logf("err2: %v", err2)
	if err2 != nil {
		t.Fatalf("err2 should be nil")
	}
}
