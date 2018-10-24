package rpc

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"math/rand"
	"reflect"
	"time"
)

const (
	basestr  = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lenReqID = 8
)

var (
	lenOfBasestr = len(basestr)
)

func strToMd5Hex(bs []byte) string {
	m := md5.New()
	m.Write(bs)
	return hex.EncodeToString(m.Sum(nil))
}

func randStr(length int) []byte {
	bs := []byte(basestr)
	result := make([]byte, 0, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bs[r.Intn(lenOfBasestr)])
	}
	return result
}

// rand request id(string) to send with Request
func randID() string {
	s := randStr(lenReqID)
	return strToMd5Hex(s)
}

// convert "in" to "out"
// "in" json.Mashal then json.Unmashal to "out"
func convert(in interface{}, out interface{}) {
	typ := reflect.TypeOf(out)
	if typ.Kind() != reflect.Ptr {
		panic("convert function err: not ptr type of out, " + typ.Kind().String())
	}

	bs, err := json.Marshal(in)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(bs, out); err != nil {
		panic(err)
	}
}

func checkArrayOrSliceInterface(v interface{}) error {
	return nil
}
