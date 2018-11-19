package rpc

import (
	"testing"
)

func Test_Server(t *testing.T) {
	s := NewServer()
	s.HandleTCP("127.0.0.1:9999")
}
