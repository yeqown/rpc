package json2

import (
	"reflect"
	"testing"
)

func Test_stdjsonCodec(t *testing.T) {
	type ST struct {
		I  int
		B2 bool
		S  string
	}
	var (
		a     = &ST{I: 10, S: "12345"}
		aPtr  = new(ST)
		codec = NewStdJSONCodec()
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
