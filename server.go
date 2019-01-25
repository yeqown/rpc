package rpc

import (
	"bufio"
	"errors"
	"log"
	"net"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/yeqown/rpc/utils"
)

var (
	typeOfError  = reflect.TypeOf((*error)(nil)).Elem()
	errEmptyBody = errors.New("empty json body")
)

// NewServerWithCodec generate a server to handle all
// tcp request from rpc client, if codec is nil will use default gobCodec
func NewServerWithCodec(codec Codec) *Server {
	if codec == nil {
		codec = newGobCodec()
	}
	return &Server{codec: codec}
}

// Server data struct to serve RPC request over TCP and HTTP
type Server struct {
	m     sync.Map // map[string]*service
	codec Codec    // codec to read request and writeResponse
}

// Register parse register type and method
// maybe save into a Map, input value is a varible
// want to got varible type name, and all Method Name
func (s *Server) Register(rcvr interface{}) error {
	srvic := new(service)
	srvic.typ = reflect.TypeOf(rcvr)
	srvic.rcvr = reflect.ValueOf(rcvr)
	sname := reflect.Indirect(srvic.rcvr).Type().Name()

	if sname == "" {
		errmsg := "rpc.Register: no service name for type " + srvic.typ.String()
		log.Print(errmsg)
		return errors.New(errmsg)
	}

	if !isExported(sname) {
		errmsg := "rpc.Register: type " + sname + " is not exported"
		log.Print(errmsg)
		return errors.New(errmsg)
	}
	srvic.name = sname
	srvic.method = suitableMethods(srvic.typ)

	if _, dup := s.m.LoadOrStore(sname, srvic); dup {
		return errors.New("rpc: service already defined: " + sname)
	}
	return nil
}

// RegisterName ... only want to export one method of rcvr
func (s *Server) RegisterName(rcvr interface{}, methodName string) error {
	srvic := new(service)
	srvic.typ = reflect.TypeOf(rcvr)
	srvic.rcvr = reflect.ValueOf(rcvr)
	sname := reflect.Indirect(srvic.rcvr).Type().Name()

	mt := suitableMethodWtihName(srvic.typ, methodName)

	i, ex := s.m.Load(sname)
	if ex {
		loadedSrvic := i.(*service)
		loadedSrvic.method[mt.method.Name] = mt
		s.m.Store(sname, loadedSrvic)
	} else {
		if sname == "" {
			errmsg := "rpc.Register: no service name for type " + srvic.typ.String()
			log.Print(errmsg)
			return errors.New(errmsg)
		}

		if !isExported(sname) {
			errmsg := "rpc.Register: type " + sname + " is not exported"
			log.Print(errmsg)
			return errors.New(errmsg)
		}
		srvic.name = sname
		srvic.method = make(map[string]*methodType)
		srvic.method[mt.method.Name] = mt
		s.m.Store(sname, srvic)
	}
	return nil
}

// before Call must parse and decode param into reflect.Value
// after Call must encode and response
func (s *Server) call(req Request) Response {
	defer func() { debugF("server called end") }()
	reqMethod := req.Method()

	dot := strings.LastIndex(reqMethod, ".") // split req.Method like "type.Method"
	if dot < 0 {
		return &defaultResponse{Err: "rpc: service/method request ill-formed: " + reqMethod, Errcode: InvalidRequest}
	}

	serviceName := reqMethod[:dot]
	methodName := reqMethod[dot+1:]

	// method existed or not
	svci, ok := s.m.Load(serviceName)
	if !ok {
		return &defaultResponse{Err: "rpc: can't find service " + reqMethod, Errcode: MethodNotFound}
	}
	svc := svci.(*service)
	mtype := svc.method[methodName]
	if mtype == nil {
		return &defaultResponse{Err: "rpc: can't find method " + reqMethod, Errcode: MethodNotFound}
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

	if err := s.codec.Decode(req.Params(), argv.Interface()); err != nil {
		debugF("decode params err: %v", err)
		return &defaultResponse{Err: err.Error(), Errcode: InvalidParamErr}
	}
	// convert(req.Params, argv.Interface())
	// fmt.Println(argv.Interface())

	replyv := reflect.New(mtype.ReplyType.Elem())
	switch mtype.ReplyType.Elem().Kind() {
	case reflect.Map:
		replyv.Elem().Set(reflect.MakeMap(mtype.ReplyType.Elem()))
	case reflect.Slice:
		replyv.Elem().Set(reflect.MakeSlice(mtype.ReplyType.Elem(), 0, 0))
	}

	if err := svc.call(mtype, argv, replyv); err != nil {
		return &defaultResponse{Err: err.Error(), Errcode: InternalErr}
	}

	byts, err := s.codec.Encode(replyv.Interface())
	if err != nil {
		return &defaultResponse{Err: err.Error(), Errcode: InternalErr}
	}
	return &defaultResponse{Err: "", Rply: byts, Errcode: SUCCESS}
}

// handleConn to recive a conn,
// parse Request and then transfer to call.
func (s *Server) handleConn(conn net.Conn) {
	// receive a request
	data, err := bufio.NewReader(conn).ReadBytes('\n')
	// debugF("recv a new request: %v", data)

	if err != nil {
		debugF("response to client connection err: %v", err)
		resp := s.codec.Response(nil, nil, InternalErr)
		utils.WriteServerTCP(conn, encodeResponse(s.codec, resp))
		return
	}

	req, err := s.codec.ParseRequest(data)
	if err := s.codec.Decode(data, req); err != nil {
		debugF("server decode request err: %v", err)
		resp := s.codec.Response(nil, nil, InvalidParamErr)
		utils.WriteServerTCP(conn, encodeResponse(s.codec, resp))
		return
	}
	debugF("[TCP] recv a new request: %s", req)

	// hanlde multi request
	if req.CanIter() {
		req.Iter(func(req Request) {
			r := s.call(req)
			r2 := s.codec.Response(req, r.Reply(), r.ErrCode())
			utils.WriteServerTCP(conn, encodeResponse(s.codec, r2))
		})
		return
	}

	r := s.call(req)
	r2 := s.codec.Response(req, r.Reply(), r.ErrCode())
	utils.WriteServerTCP(conn, encodeResponse(s.codec, r2))
}

// Start open tcp and http to serve request,
// open or not depends on that is addr an empty string,
// you can also open tcp by s.ServeTCP(addr), at the same time
// open http by s.ListenAndServe(addr)
func (s *Server) Start(tcpAddr, httpAddr string) {
	wait := make(chan bool)
	if httpAddr != "" {
		go s.ListenAndServe(httpAddr)
	}

	if tcpAddr != "" {
		go s.ServeTCP(tcpAddr)
	}

	<-wait
}

// ServeTCP Dealing with request
// decode and Call and response
func (s *Server) ServeTCP(addr string) {
	debugF("RPC server over TCP is listening: %s", addr)

	// make a listener over TCP
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
			continue
		}
		// TODO: pool goroutine
		go s.handleConn(conn)
	}
}

// ListenAndServe open http support can serve http request
func (s *Server) ListenAndServe(addr string) {
	debugF("RPC server over HTTP is listening: %s", addr)

	// TODO: replace timeout to s.codec.Response(timeoutErr).String()
	timeoutHdl := http.TimeoutHandler(s, 5*time.Second, "timeout")

	if err := http.ListenAndServe(addr, timeoutHdl); err != nil {
		panic(err)
	}
}

// ServeHTTP handle request over HTTP,
// it also implement the interface of http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		v := recover()
		if err, ok := v.(error); ok && err != nil {
			log.Printf("%v", err)
		}
	}()

	switch req.Method {
	case http.MethodGet, http.MethodPost:
		req.ParseForm()
	default:
		resp := s.codec.Response(nil, nil, MethodNotFound)
		utils.ResponseHTTP(w, encodeResponse(s.codec, resp), isDebug)
		return
	}

	data := req.Form.Get("data")
	if len(data) == 0 {
		resp := s.codec.Response(nil, nil, InvalidParamErr)
		utils.ResponseHTTP(w, encodeResponse(s.codec, resp), isDebug)
		return
	}
	debugF("[HTTP] got request data: %s", data)

	rpcReq, err := s.codec.ParseRequest([]byte(data))
	if err != nil {
		resp := s.codec.Response(nil, nil, ParseErr)
		utils.ResponseHTTP(w, encodeResponse(s.codec, resp), isDebug)
		return
	}

	// TODO: support mulit request
	resp := s.call(rpcReq)
	utils.ResponseHTTP(w, encodeResponse(s.codec, resp), isDebug)
	return
}

// encodeResponse ...
func encodeResponse(codec Codec, resp Response) []byte {
	byts, err := codec.Encode(resp)
	if err != nil {
		panic(resp)
	}

	return byts
}
