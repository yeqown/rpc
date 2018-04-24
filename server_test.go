package rpc

import (
	"testing"
)

func Test_Server(t *testing.T) {
	s := NewServer()
	s.HandleHTTP("127.0.0.1:9999")
}
