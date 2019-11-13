package rpc

import (
	"testing"
)

func Benchmark_OverTCP_jsonCodec(b *testing.B) {
	isDebug = false
	clientCodec := NewGobCodec()
	client := NewClientWithCodec(clientCodec, "127.0.0.1:9998")

	reply := 0
	args := &Args{A: 10022, B: 99999}
	for i := 0; i < b.N; i++ {
		if err := client.Call("Int.Add", args, &reply); err != nil {
			b.Errorf("call done, err=%v\n", err)
		}
	}
}

func BenchmarkOverHTTP_ServerSide(b *testing.B) {
	// isDebug = false
	// serverCodec := NewGobCodec()
	// client := NewServerWithCodec(serverCodec)

	// reply := 0
	// args := &Args{A: 10022, B: 99999}
	// for i := 0; i < b.N; i++ {
	// 	if err := client.Call("Int.Add", args, &reply); err != nil {
	// 		b.Errorf("call done, err=%v\n", err)
	// 	}
	// }
}

func BenchmarkOverHTTP_ServerSide_jsonCodec(b *testing.B) {

}
