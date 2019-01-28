package utils

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

const (
	basestr  = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lenReqID = 8
)

var l = len(basestr)

// StringToMd5Hex ...
func StringToMd5Hex(s string) string {
	m := md5.New()
	m.Write([]byte(s))
	return hex.EncodeToString(m.Sum(nil))
}

// RandString ...
func RandString(length int) string {
	bs := []byte(basestr)
	result := make([]byte, 0, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bs[r.Intn(l)])
	}
	return string(result)
}

// RandID request id(string) to send with Request
func RandID() string {
	s := RandString(lenReqID)
	return StringToMd5Hex(s)
}
