package rpc

import "testing"

// import (
// 	"reflect"
// 	"testing"
// )

// type clientTestArgs struct {
// 	A int
// 	B int
// }

// type clientTestReply struct {
// 	Sum int
// }

// func changeRPCRequest(rpcReqs []*RequestConfig) {
// 	for _, req := range rpcReqs {
// 		args := req.Args.(*clientTestArgs)
// 		reply := req.Reply.(*clientTestReply)
// 		reply.Sum = args.A + args.B
// 	}
// }

// func Test_rpcRequest(t *testing.T) {
// 	reqs := []*RequestConfig{
// 		&RequestConfig{"typ.method", &clientTestArgs{10, 10}, &clientTestReply{}},
// 		&RequestConfig{"typ.method", &clientTestArgs{22, 10}, &clientTestReply{}},
// 		&RequestConfig{"typ.method", &clientTestArgs{1212, 3434}, &clientTestReply{}},
// 	}

// 	// change in another functions
// 	changeRPCRequest(reqs)

// 	for _, req := range reqs {
// 		args := req.Args.(*clientTestArgs)
// 		sum := args.A + args.B
// 		reply := req.Reply.(*clientTestReply)

// 		t.Logf("got %v, want: %v", reply.Sum, sum)
// 		if !reflect.DeepEqual(reply, &clientTestReply{Sum: sum}) {
// 			t.Logf("got %v, want: %v", reply.Sum, sum)
// 			t.FailNow()
// 		}
// 	}
// }

func Test_parseFromRPCMethod(t *testing.T) {
	type args struct {
		reqMethod string
	}
	tests := []struct {
		name            string
		args            args
		wantServiceName string
		wantMethodName  string
		wantErr         bool
	}{
		{
			name: "case 0",
			args: args{
				reqMethod: "Int.Sum",
			},
			wantServiceName: "Int",
			wantMethodName:  "Sum",
			wantErr:         false,
		},
		{
			name: "case 1",
			args: args{
				reqMethod: "IntSum",
			},
			wantServiceName: "",
			wantMethodName:  "",
			wantErr:         true,
		},
		{
			name: "case 2",
			args: args{
				reqMethod: "Int.Sum.",
			},
			wantServiceName: "",
			wantMethodName:  "",
			wantErr:         true,
		},
		{
			name: "case 3",
			args: args{
				reqMethod: ".Int.Sum",
			},
			wantServiceName: "",
			wantMethodName:  "",
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotServiceName, gotMethodName, err := parseFromRPCMethod(tt.args.reqMethod)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseFromRPCMethod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotServiceName != tt.wantServiceName {
				t.Errorf("parseFromRPCMethod() gotServiceName = %v, want %v", gotServiceName, tt.wantServiceName)
			}
			if gotMethodName != tt.wantMethodName {
				t.Errorf("parseFromRPCMethod() gotMethodName = %v, want %v", gotMethodName, tt.wantMethodName)
			}
		})
	}
}
