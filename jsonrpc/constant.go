package jsonrpc

import (
	"github.com/yeqown/rpc"

	md5lib "crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

const (
	// VERSIONCODE const version code of JSONRPC
	VERSIONCODE = "2.0"
	basestr     = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lenReqID    = 8
)

var errcodeMap = map[int]*jsonError{
	rpc.ParseErr:        &jsonError{rpc.ParseErr, "ParseErr", nil},
	rpc.InvalidRequest:  &jsonError{rpc.InvalidRequest, "InvalidRequest", nil},
	rpc.MethodNotFound:  &jsonError{rpc.MethodNotFound, "MethodNotFound", nil},
	rpc.InvalidParamErr: &jsonError{rpc.InvalidParamErr, "InvalidParamErr", nil},
	rpc.InternalErr:     &jsonError{rpc.InternalErr, "InternalErr", nil},
}

var l = len(basestr)

// md5 ...
func md5(s string) string {
	m := md5lib.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}

// randstring ...
func randstring(length int) string {
	bs := []byte(basestr)
	result := make([]byte, 0, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bs[r.Intn(l)])
	}
	return string(result)
}

// randid request id(string) to send with NewRequest
func randid() string {
	s := randstring(lenReqID)
	return md5(s)
}
