package json2

import (
	"github.com/yeqown/rpc"
)

const (
	// VERSIONCODE const version code of JSONRPC
	VERSIONCODE = "2.0"
)

var errcodeMap = map[int]*jsonError{
	rpc.ParseErr:        &jsonError{rpc.ParseErr, "ParseErr", nil},
	rpc.InvalidRequest:  &jsonError{rpc.InvalidRequest, "InvalidRequest", nil},
	rpc.MethodNotFound:  &jsonError{rpc.MethodNotFound, "MethodNotFound", nil},
	rpc.InvalidParamErr: &jsonError{rpc.InvalidParamErr, "InvalidParamErr", nil},
	rpc.InternalErr:     &jsonError{rpc.InternalErr, "InternalErr", nil},
}
