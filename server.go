//
// 1. register function into service map
// 2. running as tcp server
// 3. accept client request, parse Args and call related function
// 4. response

package rpc

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

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

// want to save 'Type.Method' as key,
// `Method(Func type)` as value
// type MethodMap map[string]*service

func NewServer() *Server {
	return &Server{}
}

type Server struct {
	m sync.Map // map[string]*service
}

// Parse register type and method
// maybe save into a Map, input value is a varible
// want to got varible type name, and all Method Name
func (s *Server) Register(rcvr interface{}) error {

	_service := new(service)
	_service.typ = reflect.TypeOf(rcvr)
	_service.rcvr = reflect.ValueOf(rcvr)
	sname := reflect.Indirect(_service.rcvr).Type().Name()

	if sname == "" {
		err_s := "rpc.Register: no service name for type " + _service.typ.String()
		log.Print(err_s)
		return errors.New(err_s)
	}

	if !isExported(sname) {
		err_s := "rpc.Register: type " + sname + " is not exported"
		log.Print(err_s)
		return errors.New(err_s)
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
	// receive
	data, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("get an error:", err.Error())
		return
	}

	// parse request, must support multi request
	reqs := parseRequest(data)
	resps := make([]*Response, 0, MaxMultiRequest)

	// call method
	if len(reqs) > 1 {
		for _, req := range reqs {
			resp := s.call(req)
			resps = append(resps, resp)
		}
	} else {
		// single req
		req := reqs[0]
		resp := s.call(req)
		resps = append(resps, resp)
	}

	// println("len of resp: ", len(resps))
	// response to clien
	var resps_bs []byte
	if len(resps) > 1 {
		resps_bs = encodeMultiResponse(resps)
	} else {
		resps_bs = encodeResponse(resps[0])
	}
	println("response:", string(resps_bs))
	resps_bs = append(resps_bs, byte('\n'))
	conn.Write(resps_bs)
}

// Dealing with request
// decode and Call and response
func (s *Server) HandleTCP(addr string) {
	fmt.Println("start listening")
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		go s.handleConn(conn)
	}
}
