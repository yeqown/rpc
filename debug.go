package rpc

import (
	"crypto/md5"
	"hash"
	"log"
)

var (
	isDebug = true
	h       hash.Hash
)

func init() {
	h = md5.New()
}

// DebugF .
func DebugF(format string, argvs ...interface{}) {
	if !isDebug {
		return
	}
	log.Printf("[debug]: "+format, argvs...)
}

func debugHash(data []byte) string {
	if !isDebug {
		return ""
	}
	h.Reset()
	h.Write(data)
	byts := h.Sum(nil)
	return string(byts)
}
