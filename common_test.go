package rpc

import (
	"reflect"
	"testing"
)

type clientTestArgs struct {
	A int
	B int
}

type clientTestReply struct {
	Sum int
}

func changeRPCRequest(rpcReqs []*RequestConfig) {
	for _, req := range rpcReqs {
		args := req.Args.(*clientTestArgs)
		reply := req.Reply.(*clientTestReply)
		reply.Sum = args.A + args.B
	}
}

func Test_rpcRequest(t *testing.T) {
	reqs := []*RequestConfig{
		&RequestConfig{"typ.method", &clientTestArgs{10, 10}, &clientTestReply{}},
		&RequestConfig{"typ.method", &clientTestArgs{22, 10}, &clientTestReply{}},
		&RequestConfig{"typ.method", &clientTestArgs{1212, 3434}, &clientTestReply{}},
	}

	// change in another functions
	changeRPCRequest(reqs)

	for _, req := range reqs {
		args := req.Args.(*clientTestArgs)
		sum := args.A + args.B
		reply := req.Reply.(*clientTestReply)

		t.Logf("got %v, want: %v", reply.Sum, sum)
		if !reflect.DeepEqual(reply, &clientTestReply{Sum: sum}) {
			t.Logf("got %v, want: %v", reply.Sum, sum)
			t.FailNow()
		}
	}
}
