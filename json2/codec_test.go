package json2

import (
	"encoding/base64"
	"reflect"
	"testing"
)

func Test_jsonCodec(t *testing.T) {
	type ST struct {
		I  int
		B2 bool
		S  string
	}
	var (
		a     = &ST{I: 10, S: "12345"}
		aPtr  = new(ST)
		codec = NewJSONCodec()
	)

	byts, err := codec.Encode(a)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("encode success: %s\n", byts)

	if err := codec.Decode(byts, aPtr); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(a, aPtr) {
		t.Fatalf("not equal, want: %v, got: %v", a, aPtr)
	}
	t.Log("decode success\n")
}

func xTest_unmarshal(t *testing.T) {
	src := []byte("eyJJIjoxMCwiQjIiOmZhbHNlLCJTIjoiMTIzNDUifQ==")
	t.Logf("%d, %d", base64.StdEncoding.DecodedLen(len(src)), base64.StdEncoding.DecodedLen(len("eyJJIjoxMCwiQjIiOmZhbHNlLCJTIjoiMTIzNDUifQ==")))
	t.Logf("%d, %d", len(src), len("eyJJIjoxMCwiQjIiOmZhbHNlLCJTIjoiMTIzNDUifQ=="))

	dst := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
	if _, err := base64.StdEncoding.Decode(dst, src); err != nil {
		t.Fatal(err)
	}

	t.Logf("%v, %d\n", dst, len(dst))
}
