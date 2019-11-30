package rpc

import (
	"bytes"
	"reflect"
	"testing"
)

func Test_gobCodecRequest(t *testing.T) {
	codec := NewGobCodec().(*gobCodec)
	req := codec.NewRequest("svc.method", &Demo{1, "2"})
	byts, err := codec.EncodeRequests(&[]Request{req})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	reqs, err := codec.ReadRequest(byts)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !bytes.Equal(req.Params(), reqs[0].Params()) {
		t.Error("not equal params")
		t.FailNow()
	}
}

func Test_gobCodecResponse(t *testing.T) {
	codec := NewGobCodec().(*gobCodec)
	resp := codec.NewResponse(&Demo{1, "2"})
	byts, err := codec.EncodeResponses(&[]Response{resp})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	resps, err := codec.ReadResponse(byts)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !bytes.Equal(resp.Reply(), resps[0].Reply()) {
		t.Error("not equal params")
		t.FailNow()
	}
}

type Demo struct {
	VarNum int
	VarStr string
}

func Test_gobCodecEncodeDecode(t *testing.T) {
	var (
		codec      = NewGobCodec().(*gobCodec)
		bPtr       = new(bool)
		outArr     = make([]int, 0)
		outStdReqs = make([]Request, 0)
	)
	*bPtr = true

	type args struct {
		data interface{}
		out  interface{}
	}
	type test struct {
		name    string
		args    args
		wantErr bool
	}
	var tests = []test{
		{
			name: "case 0",
			args: args{
				data: &Demo{
					VarNum: 10,
					VarStr: "string",
				},
				out: &Demo{},
			},
			wantErr: false,
		},
		{
			name: "case 1",
			args: args{
				data: bPtr,
				out:  new(bool),
			},
			wantErr: false,
		},
		{
			name: "case 2",
			args: args{
				data: &[]int{1, 2, 3, 4, 5},
				out:  &outArr,
			},
			wantErr: false,
		},
		{
			name: "case 3",
			args: args{
				data: &[]Request{&stdRequest{Mthd: "a.b", Args: []byte("abc")}},
				out:  &outStdReqs,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			byts, err := codec.Encode(tt.args.data)
			if err != nil {
				t.Errorf("gobCodec.Encode() error = %v", err)
			}
			if err := codec.Decode(byts, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("gobCodec.Decode() error = %v, wantErr %v", err, tt.wantErr)
			} else if !reflect.DeepEqual(tt.args.data, tt.args.out) {
				t.Errorf("gobCodec.Decode() not equal want = %v, got %v", tt.args.data, tt.args.out)
			}
		})
	}
}

func Test_gobCodecDupEncode(t *testing.T) {
	codec := NewGobCodec().(*gobCodec)

	type ST struct {
		I int
		S string
	}
	a := &ST{
		I: 10, S: "12345",
	}

	var (
		gobEncoded, gobLastEncoded []byte
		err                        error
	)

	counter := 0
	for true {
		counter++
		if counter >= 4 {
			break
		}
		if gobEncoded, err = codec.Encode(a); err != nil {
			t.Fatal(err)
		}

		if counter == 1 {
			gobLastEncoded = gobEncoded
			continue
		}

		if !bytes.Equal(gobEncoded, gobLastEncoded) {
			t.Fatalf("dup encode 'not equal' want: %v, got: %v", gobLastEncoded, gobEncoded)
		}
		gobLastEncoded = gobEncoded
	}
	t.Log("encode dup passed")
}

func Test_gobCodecDupDecode(t *testing.T) {
	codec := NewGobCodec().(*gobCodec)

	type ST struct {
		I int
		S string
	}
	a := &ST{
		I: 10, S: "12345",
	}

	var (
		gobEncoded, _ = codec.Encode(a)
	)
	// encode
	gobEncoded, _ = codec.Encode(a)

	stPtr := new(ST)
	counter := 0

	for true {
		counter++
		if counter >= 4 {
			break
		}
		if err := codec.Decode(gobEncoded, stPtr); err != nil {
			t.Fatal(err)
		}
	}
	t.Log("decode dup passed")
}
