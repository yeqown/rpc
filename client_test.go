package rpc

import (
	"testing"
)

func Test_Client(t *testing.T) {
	s := NewServer()
	go s.HandleTCP("127.0.0.1:9999")

	c := NewClient()
	c.DialTCP("127.0.0.1:9999")
}
