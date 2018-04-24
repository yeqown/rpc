package rpc

import (
	"testing"
)

func Test_Client(t *testing.T) {
	c := NewClient()
	c.DailHTTP("127.0.0.1:9999")
}
