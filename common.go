package rpc

import (
	"errors"
	"fmt"
	"strings"
)

var (
	_ Request  = &stdRequest{}
	_ Response = &stdResponse{}
)

// Request interface contains necessary methods
type Request interface {
	// Method() to return the string of method name.
	// exmaple, "StructDemo.Method1"
	// and, should has no more form
	Method() string

	// Params means all ([]byte) data contains request args
	// all these origin params (interface{} type)
	// should be result of codec encoded
	Params() []byte

	// Ident to get indentify of request
	Ident() string
}

// Response interface contains necessary methods
type Response interface {
	// Error to return err that response struct contains
	// if there is no any error happend, should return nil
	Error() error

	// ErrorCode to return errcode(int) if ok return SUCCESS(0)
	ErrCode() int

	// Reply means all ([]byte) data contains response body
	// all these response data (interface{} type)
	// should be result of codec encoded
	Reply() []byte

	// SetReqIdent .
	SetReqIdent(ident string)
}

type stdRequest struct {
	// Mthd means Method for RPC request the method called
	Mthd string

	// Args means data to pass
	Args []byte
}

func (d *stdRequest) Method() string { return d.Mthd }
func (d *stdRequest) Params() []byte { return d.Args }
func (d *stdRequest) Ident() string  { return "" }

type stdResponse struct {
	Rply    []byte
	Err     string
	Errcode int
}

func (d *stdResponse) Error() error {
	if d.Err == "" {
		return nil
	}
	return errors.New(d.Err)
}

func (d *stdResponse) Reply() []byte            { return d.Rply }
func (d *stdResponse) ErrCode() int             { return d.Errcode }
func (d *stdResponse) SetReqIdent(ident string) { /* nothing */ }

// parseFromRPCMethod .
// split req.Method like "type.Method"
func parseFromRPCMethod(reqMethod string) (serviceName, methodName string, err error) {
	if strings.Count(reqMethod, ".") != 1 {
		return "", "", fmt.Errorf("rpc: service/method request ill-formed: %s", reqMethod)
	}
	dot := strings.LastIndex(reqMethod, ".")
	if dot < 0 {
		return "", "", fmt.Errorf("rpc: service/method request ill-formed: %s", reqMethod)
	}

	serviceName = reqMethod[:dot]
	methodName = reqMethod[dot+1:]

	return serviceName, methodName, nil
}
