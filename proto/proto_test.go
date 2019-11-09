package proto_test

import (
	"bufio"
	"bytes"
	"reflect"
	"strings"
	"testing"

	"github.com/yeqown/rpc/proto"
)

func Test_ProtoReadAndWrite(t *testing.T) {
	p := proto.New()
	p.Ver = proto.Ver2
	p.Op = proto.OpResponse
	p.Seq = uint16(22)
	p.Body = []byte("this is my body")

	buf := bytes.NewBuffer(nil)
	rr := bufio.NewReader(buf)
	wr := bufio.NewWriter(buf)

	if err := p.WriteTCP(wr); err != nil {
		t.Error(err)
		t.FailNow()
	}
	wr.Flush()
	// t.Logf("%s\n", buf.Bytes())

	p2 := proto.New()
	if err := p2.ReadTCP(rr); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(p, p2) {
		t.Errorf("not equal: want(%v) got(%v)\n", p, p2)
		t.Logf("body: want(%s), got(%s) ...\n", string(p.Body), string(p2.Body))
		t.FailNow()
	}
}

func Test_MaxSizeToTrans(t *testing.T) {
	p := proto.New()
	data := strings.Repeat("a", 1024*16)
	p.Body = []byte(data)
	println(len(p.Body))

	buf := bytes.NewBuffer(nil)
	rr := bufio.NewReader(buf)
	wr := bufio.NewWriter(buf)

	if err := p.WriteTCP(wr); err != nil {
		t.Error(err)
		t.FailNow()
	}
	wr.Flush()
	// t.Logf("%s\n", buf.Bytes())

	p2 := proto.New()
	if err := p2.ReadTCP(rr); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if !reflect.DeepEqual(p, p2) {
		// t.Errorf("not equal: want(%v) got(%v)\n", p, p2)
		t.Log(p.Ver, p.Op, p.Seq)
		t.Log(p2.Ver, p2.Op, p2.Seq)
		t.Errorf("not equal p.Body(%s), p2.Body(%s)\n", (p.Body), (p2.Body))
		t.FailNow()
	}
}
