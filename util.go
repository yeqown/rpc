package rpc

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"time"
)

const (
	BaseStr  = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ReqIdLen = 8
)

var LenOfBaseStr = len(BaseStr)

func strToMd5Hex(bs []byte) string {
	m := md5.New()
	m.Write(bs)
	return hex.EncodeToString(m.Sum(nil))
}

func randStr(length int) []byte {
	bs := []byte(BaseStr)
	result := make([]byte, 0, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bs[r.Intn(LenOfBaseStr)])
	}
	return result
}

// rand request id(string) to send with Request
func randId() string {
	s := randStr(ReqIdLen)
	return strToMd5Hex(s)
}

func convert(in interface{}, out interface{}) {
	bs, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(bs, out); err != nil {
		panic(err)
	}
}
