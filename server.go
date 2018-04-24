/*
 * mainly procedure:
 *
 * 1. register function into serviceMap
 * 2. running as http server
 * 3. accept client request, parse Args and call related function
 * 4. response
 */

package rpc

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"reflect"
)

// want to save 'Type.Method' as key, `Method(Func type)` as value
type MethodMap map[string]interface{}

func NewServer() *Server {
	return &Server{}
}

type Server struct {
	m MethodMap // method maps
}

func (s *Server) assembleKey(typ, method string) string {
	return fmt.Sprintf("%s.%s", typ, method)
}

func (s *Server) MethodMapKeys() (keys []string) {
	for k, _ := range s.m {
		keys = append(keys, k)
	}
	return
}

// Parse register type and method
// maybe save into a Map
func (s *Server) Register(value interface{}) {
	v := reflect.ValueOf(value).Elem()

	typ := v.Type().Elem()
	typName := typ.Name()

	for i := 0; i < v.NumMethod(); i++ {
		method := v.Method(i)
		methodName := method.Type().Name()
		key := s.assembleKey(typName, methodName)
		s.m[key] = method.Interface()
	}
}

/*
 * before Call must parse and decode param into reflect.Value
 * after Call must encode and response
 */
func (s *Server) Call(typMethod string, in []reflect.Value) (out []reflect.Value) {
	method, ok := s.m[typMethod]
	if !ok {
		out[0] = reflect.ValueOf("method not found")
		return
	}

	fn := reflect.ValueOf(method)
	out = fn.Call(in)

	println(len(out))
	return
}

func (s *Server) handleConn(conn io.ReadWriteCloser) {
	buf := bufio.NewWriter(conn)
	codec := &gobServerCodec{
		rwc:    conn,
		dec:    gob.NewDecoder(conn),
		enc:    gob.NewDecoder(buf),
		encBuf: buf,
	}

	// if err := codec.ReadRequest()

	// receive
	data, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("get an error:", err.Error())
		return
	}
	fmt.Println("Message Received", data)
	// response
	conn.Write([]byte(data + "\n"))
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
