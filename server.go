package rpc

import (
	"bufio"
	"errors"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
)

var (
	typeOfError  = reflect.TypeOf((*error)(nil)).Elem()
	errEmptyBody = errors.New("empty json body")
)

// NewServerWithCodec generate a server to handle all
// tcp request from rpc client, if codec is nil will use default gobCodec
func NewServerWithCodec(addr string, codec Codec) *Server {
	if codec == nil {
		codec = newGobCodec()
	}
	return &Server{
		addr:  addr,
		codec: codec,
	}
}

// Server data struct to serve RPC request over TCP and HTTP
type Server struct {
	addr  string   // addr to tcp listen on
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
		return &defaultResponse{
			Err: "rpc: service/method request ill-formed: " + reqMethod,
		}
	}

	serviceName := reqMethod[:dot]
	methodName := reqMethod[dot+1:]

	// method existed or not
	svci, ok := s.m.Load(serviceName)
	if !ok {
		return &defaultResponse{
			Err: "rpc: can't find service " + reqMethod,
		}
	}
	svc := svci.(*service)
	mtype := svc.method[methodName]
	if mtype == nil {
		return &defaultResponse{
			Err: "rpc: can't find method " + reqMethod,
		}
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
		return &defaultResponse{Err: err.Error()}
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
		return &defaultResponse{Err: err.Error()}
	}

	byts, err := s.codec.Encode(replyv.Interface())
	if err != nil {
		return &defaultResponse{Err: err.Error()}
	}
	return &defaultResponse{Err: "", Rply: byts}
}

// handleConn to recive a conn,
// parse Request and then transfer to call.
func (s *Server) handleConn(conn net.Conn) {
	// receive a request
	data, err := bufio.NewReader(conn).ReadBytes('\n')
	debugF("recv a new request: %v", data)

	if err != nil {
		debugF("response to client connection err: %v", err)
		s.codec.Response(conn, nil, nil, err)
		return
	}

	req, err := s.codec.ParseRequest(data)
	if err := s.codec.Decode(data, req); err != nil {
		debugF("server decode request err: %v", err)
		s.codec.Response(conn, nil, nil, err)
		return
	}

	// hanlde multi request
	if req.CanIter() {
		req.Iter(func(req Request) {
			resp := s.call(req)
			s.codec.Response(conn, req, resp.Reply(), resp.Error())
		})
		return
	}

	resp := s.call(req)
	s.codec.Response(conn, req, resp.Reply(), resp.Error())
}

// func (s *Server) handleWithRequests(reqs []*Request) []*Response {
// 	resps := make([]*Response, 0, MaxMultiRequest)
// 	// call method
// 	if len(reqs) > 1 {
// 		for _, req := range reqs {
// 			resp := s.call(req)
// 			resps = append(resps, resp)
// 		}
// 	} else {
// 		req := reqs[0]
// 		resp := s.call(req)
// 		resps = append(resps, resp)
// 	}
// 	return resps
// }

// ServeTCP Dealing with request
// decode and Call and response
func (s *Server) ServeTCP() {
	debugF("RPC Server is listening: %s", s.addr)

	// make a listener over TCP
	listener, err := net.Listen("tcp", s.addr)
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

// Handle Request over HTTP
// Inspired by `https://github.com/gorilla/rpc`
// func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	// valid request method
// 	var resp *Response

// 	w.Header().Set("Content-Type", "application/json")

// 	if r.Method != http.MethodPost {
// 		resp = NewResponse("", nil, NewJsonrpcErr(
// 			http.StatusMethodNotAllowed, "HTTP request method must be POST", nil),
// 		)
// 		response(w, resp)
// 		return
// 	}
// 	// parseJsonBody
// 	reqs, err := getRequestFromBody(r)
// 	if err != nil {
// 		resp = NewResponse("", nil, NewJsonrpcErr(InternalErr, err.Error(), nil))
// 		response(w, resp)
// 		return
// 	}

// 	resps := s.handleWithRequests(reqs)

// 	if len(resps) > 1 {
// 		response(w, resps)
// 	} else {
// 		response(w, resps[0])
// 	}
// 	return
// }

// response
// func response(w http.ResponseWriter, i interface{}) {
// 	bs, err := json.Marshal(i)
// 	if err != nil {
// 		resp := NewResponse("", nil,
// 			NewJsonrpcErr(InternalErr, err.Error(), nil),
// 		)
// 		bs, _ = json.Marshal(resp)
// 	}
// 	_, err = io.WriteString(w, string(bs))
// 	if err != nil {
// 		println("Send to client err:", err.Error())
// 	}
// }

// getRequestFromBody support parse request from jsonBody
// and parse into Request
// func getRequestFromBody(req *http.Request) ([]*Request, error) {
// 	var (
// 		body []byte
// 		err  error
// 	)
// 	if body, err = ioutil.ReadAll(req.Body); err != nil {
// 		return nil, err
// 	}
// 	if len(body) == 0 {
// 		return nil, errEmptyBody
// 	}
// 	// parse []byte into Request
// 	mReq, err := parseRequest(body)
// 	return mReq, err
// }
