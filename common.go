package rpc

import "errors"

var (
	_ Request  = &defaultRequest{}
	_ Response = &defaultResponse{}
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

	// CanIter tell the request is multi or not
	CanIter() bool

	// Iter is func to loop all request in special request
	Iter(func(req Request))
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
}

type defaultRequest struct {

	// Mthd means Method for RPC request the method called
	Mthd string

	// Args means data to pass
	Args []byte
}

func (d *defaultRequest) Method() string                  { return d.Mthd }
func (d *defaultRequest) Params() []byte                  { return d.Args }
func (d *defaultRequest) CanIter() bool                   { return false }
func (d *defaultRequest) Iter(iterFunc func(req Request)) { return }

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

func (d *defaultResponse) Reply() []byte { return d.Rply }
func (d *defaultResponse) ErrCode() int  { return d.Errcode }
