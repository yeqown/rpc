package rpc

import (
	"bufio"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"reflect"
	"runtime/debug"
	"sync"
	"time"

	"github.com/yeqown/rpc/proto"
)

var (
	typeOfError  = reflect.TypeOf((*error)(nil)).Elem()
	errEmptyBody = errors.New("empty json body")
)

// NewServerWithCodec generate a server to handle all
// tcp request from rpc client, if codec is nil will use default gobCodec
func NewServerWithCodec(codec ServerCodec) *Server {
	if codec == nil {
		codec = NewGobCodec()
	}
	return &Server{codec: codec}
}

// Server data struct to serve RPC request over TCP and HTTP
type Server struct {
	m     sync.Map    // map[string]*service
	codec ServerCodec // codec to read request and writeResponse
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
		// log.Print(errmsg)
		DebugF(errmsg)
		return errors.New(errmsg)
	}

	if !isExported(sname) {
		errmsg := "rpc.Register: type " + sname + " is not exported"
		// log.Print(errmsg)
		DebugF(errmsg)
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
			DebugF(errmsg)
			return errors.New(errmsg)
		}

		if !isExported(sname) {
			errmsg := "rpc.Register: type " + sname + " is not exported"
			DebugF(errmsg)
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
func (s *Server) call(reqs []Request) (replies []Response) {
	defer func() { DebugF("server called end") }()
	replies = make([]Response, len(reqs))
	for idx, req := range reqs {
		var (
			reply Response
		)

		serviceName, methodName, err := parseFromRPCMethod(req.Method())
		if err != nil {
			DebugF("parseFromRPCMethod err=%v", err)
			reply = s.codec.ErrResponse(InvalidRequest, err)
			replies[idx] = reply
			continue
			// goto errSkip
		}

		// method existed or not
		svci, ok := s.m.Load(serviceName)
		if !ok {
			err := errors.New("rpc: can't find service " + serviceName)
			reply = s.codec.ErrResponse(MethodNotFound, err)
			replies[idx] = reply
			continue
		}

		svc := svci.(*service)
		mtype := svc.method[methodName]
		if mtype == nil {
			err := errors.New("rpc: can't find method " + req.Method())
			reply = s.codec.ErrResponse(MethodNotFound, err)
			replies[idx] = reply
			continue
		}

		// To prepare argv and replyv in reflect.Value. Refer to `net/http/rpc`
		// If true, need to indirect before calling.
		var (
			argv       reflect.Value
			argIsValue = false
		)
		if mtype.ArgType.Kind() == reflect.Ptr {
			argv = reflect.New(mtype.ArgType.Elem())
		} else {
			argv = reflect.New(mtype.ArgType)
			argIsValue = true
		}
		if argIsValue {
			argv = argv.Elem() // argv guaranteed to be a pointer now.
		}

		if err := s.codec.ReadRequestBody(req.Params(), argv.Interface()); err != nil {
			DebugF("could not readRequestBody err=%v", err)
			err := errors.New("rpc: could not read request body " + req.Method())
			reply = s.codec.ErrResponse(InternalErr, err)
			replies[idx] = reply
			continue
		}

		var replyv reflect.Value
		replyv = reflect.New(mtype.ReplyType.Elem())
		switch mtype.ReplyType.Elem().Kind() {
		case reflect.Map:
			replyv.Elem().Set(reflect.MakeMap(mtype.ReplyType.Elem()))
		case reflect.Slice:
			replyv.Elem().Set(reflect.MakeSlice(mtype.ReplyType.Elem(), 0, 0))
		}

		if err := svc.call(mtype, argv, replyv); err != nil {
			reply = s.codec.ErrResponse(InternalErr, err)
		} else {
			// normal response
			reply = s.codec.NewResponse(replyv.Interface())
			reply.SetReqIdent(req.Ident())
		}

		replies[idx] = reply
	}
	// for req range reqs. END

	return
}

// serveConn to recive a conn,
// parse NewRequest and then transfer to call.
func (s *Server) serveConn(conn net.Conn) {
	// receive a request
	// data, err := bufio.NewReader(conn).ReadBytes('\n')
	rr := bufio.NewReader(conn)
	wr := bufio.NewWriter(conn)

	var (
		precv = proto.New()
		psend = proto.New()
		resps = make([]Response, 0)
	)

	for {
		if err := precv.ReadTCP(rr); err != nil {
			// true: tcp closed
			DebugF("ReadTCP error: %v", err)
			break
		}

		DebugF("recv a new request: %v", precv.Body)
		reqs, err := s.codec.ReadRequest(precv.Body)
		if err != nil {
			DebugF("could not parse request: %v", err)
			resps = append(resps, s.codec.ErrResponse(ParseErr, err))
			goto wr
		}
		resps = s.call(reqs)
		DebugF("s.call(req) req: %v result: %v", reqs, resps)
	wr:
		if psend.Body, err = s.codec.EncodeResponses(resps); err != nil {
			DebugF("could not encode responses, err=%v", err)
			continue
		}
		psend.WriteTCP(wr)
		wr.Flush()
	}
}

// ServeTCP Dealing with request
// decode and Call and response
func (s *Server) ServeTCP(addr string) {
	DebugF("RPC server over TCP is listening: %s", addr)

	// make a listener over TCP
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			DebugF("listener.Accept(), err=%v", err)
			continue
		}
		// TODO: goroutine pool
		go s.serveConn(conn)
	}
}

// ListenAndServe open http support can serve http request
func (s *Server) ListenAndServe(addr string) {
	log.Printf("RPC server over HTTP is listening: %s", addr)
	if err := http.ListenAndServe(
		addr,
		http.TimeoutHandler(s, 5*time.Second, "timeout"),
	); err != nil {
		panic(err)
	}
}

// ServeHTTP handle request over HTTP,
// it also implement the interface of http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err, ok := recover().(error); ok && err != nil {
			// log.Printf("%v", err)
			DebugF("[ServeHTTP] recover %v with stack: \n", err)
			debug.PrintStack()
		}
	}()

	var (
		data []byte
		err  error
	)
	switch req.Method {
	case http.MethodPost:
		if data, err = ioutil.ReadAll(req.Body); err != nil {
			resp := s.codec.ErrResponse(InvalidParamErr, err)
			byts, err := s.codec.EncodeResponses(resp)
			DebugF("s.codec.EncodeResponses err=%v", err)
			String(w, http.StatusOK, byts)
			// JSON(w, http.StatusOK, resp)
			return
		}
		defer req.Body.Close()
	default:
		err := errors.New("method not allowed: " + req.Method)
		resp := s.codec.ErrResponse(MethodNotFound, err)
		// JSON(w, http.StatusOK, resp)
		byts, err := s.codec.EncodeResponses(resp)
		DebugF("s.codec.EncodeResponses err=%v", err)
		String(w, http.StatusOK, byts)
		return
	}

	DebugF("[HTTP] got request data: %v", data)
	rpcReqs, err := s.codec.ReadRequest([]byte(data))
	if err != nil {
		resp := s.codec.ErrResponse(ParseErr, err)
		// JSON(w, http.StatusOK, resp)
		byts, err := s.codec.EncodeResponses(resp)
		DebugF("s.codec.EncodeResponses err=%v", err)
		String(w, http.StatusOK, byts)
		return
	}

	resps := s.call(rpcReqs)
	DebugF("s.call(rpcReq) result: %v", resps)
	// JSON(w, http.StatusOK, resps)
	byts, err := s.codec.EncodeResponses(resps)
	DebugF("s.codec.EncodeResponses err=%v", err)
	String(w, http.StatusOK, byts)
	return
}

// JSON .
func JSON(w http.ResponseWriter, statusCode int, v interface{}) error {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	byts, err := json.Marshal(v)
	if err != nil {
		DebugF("could not marshal v=%v, err=%v", v, err)
		return err
	}

	_, err = io.WriteString(w, string(byts))
	return err
}

// String .
func String(w http.ResponseWriter, statusCode int, byts []byte) error {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "text/plain")
	// byts, err := json.Marshal(v)
	// if err != nil {
	// 	DebugF("could not marshal v=%v, err=%v", v, err)
	// 	return err
	// }

	_, err := io.WriteString(w, string(byts))
	return err
}
