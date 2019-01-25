package json2

import (
	"bytes"
	"encoding/json"
	"testing"
)

type TestArgs struct {
	TypeInt    int
	TypeString string
	// Next       *TestArgs
}

var testByts []byte

func init() {
	args := &TestArgs{
		TypeInt:    10,
		TypeString: "string",
		// Next:       &TestArgs{11, "string11", nil},
	}
	testByts, _ = json.Marshal(args)
}

func Test_jsonRequest(t *testing.T) {
	req := &jsonRequest{
		ID:      "id11",
		Mthd:    "Typ.Mthd",
		Args:    testByts,
		Version: VERSIONCODE,
	}

	if req.Next() != nil {
		t.Log("must has no next")
		t.Fail()
	}

	if req.HasNext() {
		t.Log("must has no next")
		t.Fail()
	}
}

func Test_jsonResponse(t *testing.T) {

}

func Test_jsonRequetArray(t *testing.T) {
	arr := []*jsonRequest{
		&jsonRequest{"id10", "Typ.Mthd", testByts, VERSIONCODE},
		&jsonRequest{"id11", "Typ.Mthd", testByts, VERSIONCODE},
		&jsonRequest{"id12", "Typ.Mthd", testByts, VERSIONCODE},
		&jsonRequest{"id13", "Typ.Mthd", testByts, VERSIONCODE},
		&jsonRequest{"id14", "Typ.Mthd", testByts, VERSIONCODE},
	}

	req := &jsonRequestArray{
		cur:  0,
		Data: arr,
	}

	if !req.HasNext() {
		t.Log("req has next!!!")
		t.FailNow()
	}

	counter := 0
	for req.HasNext() {
		counter++
		v := req.Next()
		jreq, ok := v.(*jsonRequest)
		t.Logf("jsonRequest iter %s\n", jreq)
		if !ok {
			t.Log("nani yo?")
			t.FailNow()
		}

		if !bytes.Equal(jreq.Args, testByts) {
			t.Log("not equal ?")
			t.FailNow()
		}
	}

	if counter != len(req.Data) {
		t.Logf("iter not finish: got %d, want: %d\n", counter, len(req.Data))
		t.FailNow()
	}
}

func Test_jsonResponseArray(t *testing.T) {
	arr := []*jsonResponse{
		&jsonResponse{"id10", nil, testByts, VERSIONCODE},
		&jsonResponse{"id11", nil, testByts, VERSIONCODE},
		&jsonResponse{"id12", nil, testByts, VERSIONCODE},
		&jsonResponse{"id13", nil, testByts, VERSIONCODE},
		&jsonResponse{"id14", nil, testByts, VERSIONCODE},
	}

	req := &jsonResponseArray{
		cur:  0,
		Data: arr,
	}

	if !req.HasNext() {
		t.Log("req has next!!!")
		t.FailNow()
	}

	counter := 0
	for req.HasNext() {
		counter++
		v := req.Next()
		jreq, ok := v.(*jsonResponse)
		t.Logf("jsonResponse iter %s\n", jreq)
		if !ok {
			t.Log("nani yo?")
			t.FailNow()
		}

		if !bytes.Equal(jreq.Result, testByts) {
			t.Log("not equal ?")
			t.FailNow()
		}
	}

	if counter != len(req.Data) {
		t.Logf("iter not finish: got %d, want: %d\n", counter, len(req.Data))
		t.FailNow()
	}
}
