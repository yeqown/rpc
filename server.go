//
// 1. register function into service map
// 2. running as tcp server
// 3. accept client request, parse Args and call related function
// 4. response

package rpc

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

var (
	typeOfError  = reflect.TypeOf((*error)(nil)).Elem()
	errEmptyBody = errors.New("empty json body")
)

type methodType struct {
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
	// sync.Mutex // protects counters
	// numCalls   uint
}

type service struct {
	name   string
	rcvr   reflect.Value
	typ    reflect.Type
	method map[string]*methodType
}

func (s *service) call(mtype *methodType, req *Request, argv, replyv reflect.Value) *Response {
	function := mtype.method.Func
	// fmt.Println(argv, replyv)
	returnValues := function.Call([]reflect.Value{s.rcvr, argv, replyv})
	errIter := returnValues[0].Interface()

	errmsg := ""
	if errIter != nil {
		errmsg = errIter.(error).Error()
		return NewResponse(req.ID, nil, NewJsonrpcErr(InternalErr, errmsg, nil))
	}

	return NewResponse(req.ID, replyv.Interface(), nil)
}

// NewServer want to save 'Type.Method' as key,
// `Method(Func type)` as value
// type MethodMap map[string]*service
func NewServer() *Server {
	return &Server{}
}

// Server ...
type Server struct {
	m sync.Map // map[string]*service
}

// Register parse register type and method
// maybe save into a Map, input value is a varible
// want to got varible type name, and all Method Name
func (s *Server) Register(rcvr interface{}) error {

	_service := new(service)
	_service.typ = reflect.TypeOf(rcvr)
	_service.rcvr = reflect.ValueOf(rcvr)
	sname := reflect.Indirect(_service.rcvr).Type().Name()

	if sname == "" {
		errmsg := "rpc.Register: no service name for type " + _service.typ.String()
		log.Print(errmsg)
		return errors.New(errmsg)
	}

	if !isExported(sname) {
		errmsg := "rpc.Register: type " + sname + " is not exported"
		log.Print(errmsg)
		return errors.New(errmsg)
	}
	_service.name = sname
	_service.method = suitableMethods(_service.typ, true)

	if _, dup := s.m.LoadOrStore(sname, _service); dup {
		return errors.New("rpc: service already defined: " + sname)
	}
	return nil
}

// suitableMethods get all method of registering-type
// into a map[string]*methodType
func suitableMethods(typ reflect.Type, reportErr bool) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name

		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// Method needs three ins: receiver, *args, *reply.
		if mtype.NumIn() != 3 {
			if reportErr {
				log.Printf("rpc.Register: method %q has %d input parameters; needs exactly three\n", mname, mtype.NumIn())
			}
			continue
		}
		// First arg need not be a pointer.
		argType := mtype.In(1)
		if !isExportedOrBuiltinType(argType) {
			if reportErr {
				log.Printf("rpc.Register: argument type of method %q is not exported: %q\n", mname, argType)
			}
			continue
		}
		// Second arg must be a pointer.
		replyType := mtype.In(2)
		if replyType.Kind() != reflect.Ptr {
			if reportErr {
				log.Printf("rpc.Register: reply type of method %q is not a pointer: %q\n", mname, replyType)
			}
			continue
		}
		// Reply type must be exported.
		if !isExportedOrBuiltinType(replyType) {
			if reportErr {
				log.Printf("rpc.Register: reply type of method %q is not exported: %q\n", mname, replyType)
			}
			continue
		}
		// Method needs one out.
		if mtype.NumOut() != 1 {
			if reportErr {
				log.Printf("rpc.Register: method %q has %d output parameters; needs exactly one\n", mname, mtype.NumOut())
			}
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != typeOfError {
			if reportErr {
				log.Printf("rpc.Register: return type of method %q is %q, must be error\n", mname, returnType)
			}
			continue
		}
		methods[mname] = &methodType{method: method, ArgType: argType, ReplyType: replyType}
	}
	return methods
}

// type isExported
func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

// Is this type exported or a builtin?
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}

// before Call must parse and decode param into reflect.Value
// after Call must encode and response
func (s *Server) call(req *Request) *Response {
	// TODO: simplfy this function, or split into several functions
	dot := strings.LastIndex(req.Method, ".") // split req.Method like "type.Method"
	if dot < 0 {
		err := errors.New("rpc: service/method request ill-formed: " + req.Method)
		return NewResponse(req.ID, nil, NewJsonrpcErr(ParseErr, err.Error(), err))
	}

	serviceName := req.Method[:dot]
	methodName := req.Method[dot+1:]

	// method existed or not
	svci, ok := s.m.Load(serviceName)
	if !ok {
		err := errors.New("rpc: can't find service " + req.Method)
		return NewResponse(req.ID, nil, NewJsonrpcErr(MethodNotFound, err.Error(), nil))
	}
	svc := svci.(*service)
	mtype := svc.method[methodName]
	if mtype == nil {
		err := errors.New("rpc: can't find method " + req.Method)
		return NewResponse(req.ID, nil, NewJsonrpcErr(MethodNotFound, err.Error(), nil))
	}

	// to prepare argv and replyv in reflect.Value
	// ref to `net/http/rpc`
	argIsValue := false // if true, need to indirect before calling.
	var argv reflect.Value
	if mtype.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(mtype.ArgType.Elem())
	} else {
		argv = reflect.New(mtype.ArgType)
		argIsValue = true
	}

	// argv guaranteed to be a pointer now.
	if argIsValue {
		argv = argv.Elem()
	}

	convert(req.Params, argv.Interface())
	// fmt.Println(argv.Interface())

	replyv := reflect.New(mtype.ReplyType.Elem())
	switch mtype.ReplyType.Elem().Kind() {
	case reflect.Map:
		replyv.Elem().Set(reflect.MakeMap(mtype.ReplyType.Elem()))
	case reflect.Slice:
		replyv.Elem().Set(reflect.MakeSlice(mtype.ReplyType.Elem(), 0, 0))
	}

	return svc.call(mtype, req, argv, replyv)
}

// handleConn to recive a conn,
// parse Request and then transfer to call.
func (s *Server) handleConn(conn io.ReadWriteCloser) {
	var respsBytes []byte

	// receive
	data, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		errmsg := "reciving connection get an error:" + err.Error()
		fmt.Println(errmsg)
		respsBytes = encodeResponse(
			NewResponse("", nil,
				NewJsonrpcErr(ParseErr, errmsg, nil),
			),
		)
		conn.Write(respsBytes)
		return
	}

	// parse request, must support multi request
	reqs, err := parseRequest(data)
	if err != nil {
		respsBytes = encodeResponse(
			NewResponse("", nil,
				NewJsonrpcErr(ParseErr, err.Error(), nil),
			),
		)
		conn.Write(respsBytes)
		return
	}

	resps := s.handleWithRequests(reqs)
	if len(resps) > 1 {
		respsBytes = encodeMultiResponse(resps)
	} else {
		respsBytes = encodeResponse(resps[0])
	}

	// println("response:", string(respsBytes))
	respsBytes = append(respsBytes, byte('\n'))
	conn.Write(respsBytes)
}

func (s *Server) handleWithRequests(reqs []*Request) []*Response {
	resps := make([]*Response, 0, MaxMultiRequest)

	// call method
	if len(reqs) > 1 {
		for _, req := range reqs {
			resp := s.call(req)
			resps = append(resps, resp)
		}
	} else {
		req := reqs[0]
		resp := s.call(req)
		resps = append(resps, resp)
	}
	return resps
}

// HandleTCP Dealing with request
// decode and Call and response
func (s *Server) HandleTCP(addr string) {
	log.Println("RPC Server is listening")
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		go s.handleConn(conn)
	}
}

// Handle Request over HTTP
// Inspired by `https://github.com/gorilla/rpc`
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// valid request method
	var resp *Response

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		resp = NewResponse("", nil, NewJsonrpcErr(
			http.StatusMethodNotAllowed, "HTTP request method must be POST", nil),
		)
		response(w, resp)
		return
	}
	// parseJsonBody
	reqs, err := getRequestFromBody(r)
	if err != nil {
		resp = NewResponse("", nil, NewJsonrpcErr(InternalErr, err.Error(), nil))
		response(w, resp)
		return
	}

	resps := s.handleWithRequests(reqs)

	if len(resps) > 1 {
		response(w, resps)
	} else {
		response(w, resps[0])
	}
	return
}

// response
func response(w http.ResponseWriter, i interface{}) {
	bs, err := json.Marshal(i)
	if err != nil {
		resp := NewResponse("", nil,
			NewJsonrpcErr(InternalErr, err.Error(), nil),
		)
		bs, _ = json.Marshal(resp)
	}
	_, err = io.WriteString(w, string(bs))
	if err != nil {
		println("Send to client err:", err.Error())
	}
}

// getRequestFromBody support parse request from jsonBody
// and parse into Request
func getRequestFromBody(req *http.Request) ([]*Request, error) {
	var (
		body []byte
		err  error
	)
	if body, err = ioutil.ReadAll(req.Body); err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, errEmptyBody
	}
	// parse []byte into Request
	mReq, err := parseRequest(body)
	return mReq, err
}
