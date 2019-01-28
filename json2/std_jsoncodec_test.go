package json2

// import (
// 	"reflect"
// 	"testing"

// 	"github.com/yeqown/rpc"
// )

// func Test_stdjsonCodec(t *testing.T) {
// 	type ST struct {
// 		I  int
// 		B2 bool
// 		S  string
// 	}
// 	var (
// 		a     = &ST{I: 10, S: "12345"}
// 		aPtr  = new(ST)
// 		codec = NewStdJSONCodec()
// 	)

// 	byts, err := codec.Encode(a)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Logf("encode success: %s\n", byts)

// 	if err := codec.Decode(byts, aPtr); err != nil {
// 		t.Fatal(err)
// 	}

// 	if !reflect.DeepEqual(a, aPtr) {
// 		t.Fatalf("not equal, want: %v, got: %v", a, aPtr)
// 	}
// 	t.Log("decode success\n")
// }

// func Test_stdjsonCodecjsonRequest(t *testing.T) {

// 	codec := NewStdJSONCodec()
// 	req := codec.Request("typ.Method", []byte("hahah"))
// 	byts, err := codec.Encode(req)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	newReq, err := codec.ParseRequest(byts)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if !reflect.DeepEqual(req, newReq) {
// 		t.Logf("codec.ParseRequest(byts) want: %v, got: %v", req, newReq)
// 	}

// 	// dup decode ...
// 	newReq2, err2 := codec.ParseRequest(byts)
// 	if err2 != nil {
// 		t.Fatal(err2)
// 	}
// 	if !reflect.DeepEqual(req, newReq2) {
// 		t.Logf("codec.ParseRequest(byts) want: %v, got: %v", req, newReq)
// 	}
// }

// func Test_stdjsonCodecRequestMulti(t *testing.T) {
// 	var codec rpc.Codec
// 	codec = NewStdJSONCodec()

// 	t.Log(codec.MultiSupported())

// 	type DemoArgs struct {
// 		A int `json:"a"`
// 		B int `json:"b"`
// 	}

// 	testByts, _ := codec.Encode(&DemoArgs{10, 13})

// 	reqs := []rpc.Request{
// 		&jsonRequest{
// 			ID:      "123810",
// 			Mthd:    "typ.Method",
// 			Args:    testByts,
// 			Version: VERSIONCODE,
// 		},
// 		&jsonRequest{
// 			ID:      "123810",
// 			Mthd:    "typ.Method",
// 			Args:    testByts,
// 			Version: VERSIONCODE,
// 		},
// 		&jsonRequest{
// 			ID:      "123810",
// 			Mthd:    "typ.Method",
// 			Args:    testByts,
// 			Version: VERSIONCODE,
// 		},
// 	}

// 	finalReq := codec.RequestMulti(reqs)
// 	testRequestArrayIter(t, finalReq, testByts)

// 	if byts, err := codec.Encode(finalReq); err != nil {
// 		t.Log("codec.Encode(finalReq)")
// 		t.FailNow()
// 	} else {
// 		var req rpc.Request
// 		if req, err = codec.ParseRequest(byts); err != nil {
// 			t.Logf("codec.ParseRequest(byts): %s\n", byts)
// 			t.FailNow()
// 		}
// 		t.Logf("parseRequest as: %v", req)
// 		testRequestArrayIter(t, req, testByts)
// 	}
// }

// func testRequestArrayIter(t *testing.T, req rpc.Request, testByts []byte) {
// 	t.Logf("testRequestArrayIter %v", req)
// 	// Iter interface test
// 	if !req.HasNext() {
// 		t.Log("request multi bug, cannot support HasNext functions")
// 		t.FailNow()
// 	}

// 	demo := &jsonRequest{
// 		ID:      "123810",
// 		Mthd:    "typ.Method",
// 		Args:    testByts,
// 		Version: VERSIONCODE,
// 	}

// 	for req.HasNext() {
// 		v := req.Next()
// 		t.Logf("req.Next() got: %v, want: %v", v, demo)
// 		// if !reflect.DeepEqual(v, demo) {
// 		// 	t.FailNow()
// 		// }
// 	}
// }

// type AArgs struct {
// 	A int `json:"a"`
// 	B int `json:"b"`
// }

// // func Test_decode(t *testing.T) {
// // 	codec := NewStdJSONCodec()
// // 	src := []byte("eyJhIjoyMTMxMiwiYiI6MTkw")
// // 	a := new(AArgs)

// // 	if err := codec.Decode(src, a); err != nil {
// // 		t.Fatal(err)
// // 	}
// // 	t.Logf("%v", a)
// // }
