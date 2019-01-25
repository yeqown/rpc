package rpc

import "errors"

var (
	_ Request  = &defaultRequest{}
	_ Response = &defaultResponse{}
)

// Iter interface means data structure can be itered
type Iter interface {
	HasNext() bool
	Next() interface{}
}

// Request interface contains necessary methods
type Request interface {

	// Iter interface contains 'HasNext() bool' and 'Next() interface{}'
	Iter

	// Method() to return the string of method name.
	// exmaple, "StructDemo.Method1"
	// and, should has no more form
	Method() string

	// Params means all ([]byte) data contains request args
	// all these origin params (interface{} type)
	// should be result of codec encoded
	Params(codec Codec) []byte
}

// Response interface contains necessary methods
type Response interface {
	// Iter interface contains 'HasNext() bool' and 'Next() interface{}'
	Iter

	// Error to return err that response struct contains
	// if there is no any error happend, should return nil
	Error() error

	// ErrorCode to return errcode(int) if ok return SUCCESS(0)
	ErrCode() int

	// Reply means all ([]byte) data contains response body
	// all these response data (interface{} type)
	// should be result of codec encoded
	Reply(codec Codec) []byte
}

type defaultRequest struct {

	// Mthd means Method for RPC request the method called
	Mthd string

	// Args means data to pass
	Args []byte
}

func (d *defaultRequest) Method() string           { return d.Mthd }
func (d *defaultRequest) Params(code Codec) []byte { return d.Args }
func (d *defaultRequest) HasNext() bool            { return false }
func (d *defaultRequest) Next() interface{}        { return nil }

type defaultResponse struct {
	Rply    []byte
	Err     string
	Errcode int
}

func (d *defaultResponse) Error() error {
	if d.Err == "" {
		return nil
	}
	return errors.New(d.Err)
}

func (d *defaultResponse) Reply(code Codec) []byte { return d.Rply }
func (d *defaultResponse) ErrCode() int            { return d.Errcode }
func (d *defaultResponse) HasNext() bool           { return false }
func (d *defaultResponse) Next() interface{}       { return nil }
