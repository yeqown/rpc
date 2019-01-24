package rpc

import (
	"errors"
	"fmt"
	"net"
)

// ArgsEncodeFunc ...
// type ArgsEncodeFunc func(args interface{}) ([]byte, error)

var (
	errMultiReplyTypePtr = errors.New("multi reply should be arrry or slice pointer")
	errEmptyCodec        = errors.New("client has an empty codec")
	argsEncodeFunc       = defaultArgsEncodeFunc
)

// TODO: finish this
func defaultArgsEncodeFunc() ([]byte, error) {
	return nil, nil
}

// NewClient ... maybe noneed this function
func NewClient(addr string) *Client {
	return &Client{
		addr:  addr,
		codec: newGobCodec(),
	}
}

// NewClientWithCodec ...
func NewClientWithCodec(addr string, c Codec) *Client {
	return &Client{
		addr:  addr,
		codec: c,
	}
}

// Client ....
type Client struct {
	addr  string
	codec Codec
	conn  net.Conn // connection to server
}

// Call Client Call Remote method
// TODO: timeout cancel
func (c *Client) Call(method string, args, reply interface{}) error {
	if c.codec == nil {
		return errEmptyCodec
	}

	// connect to server
	if c.conn == nil {
		conn, err := net.Dial("tcp", c.addr)
		if err != nil {
			return fmt.Errorf("net.Dial tcp get err: %v", err)
		}
		c.conn = conn
	}

	// core ....
	respDataByts, err := c.codec.Request(c.conn, method, args)
	if err != nil {
		return fmt.Errorf("c.codec.Request(c.conn, req) got err: %v", err)
	}

	var resp Response
	if resp, err = c.codec.ParseResponse(respDataByts); err != nil {
		return err
	}
	if err := resp.Error(); err != nil {
		return fmt.Errorf("resp.Error(): %v", err)
	}
	// core ...

	if err := c.codec.Decode(resp.Reply(), reply); err != nil {
		return fmt.Errorf("c.codec.Decode(resp.Reply() got err: %v", err)
	}
	return nil
}

// Close the client connectio to the server
func (c *Client) Close() {
	if c.conn == nil {
		return
	}
	c.conn.Close()
}

// CallMulti ...
// TODO: handle with multi params and multi response? and how
// func (c *Client) CallMulti(method string, params, replys interface{}) error {
// 	ele := reflect.ValueOf(params)
// 	typ := reflect.TypeOf(params)

// 	if typ.Kind() == reflect.Ptr {
// 		ele = ele.Elem()
// 		typ = ele.Type()
// 	}

// 	if typ.Kind() != reflect.Slice && typ.Kind() != reflect.Array {
// 		err := fmt.Errorf("Error: params type %s is not array type or slice", typ.Kind().String())
// 		return err
// 	}

// 	reqs := make([]*Request, 0)
// 	for i := 0; i < ele.Len(); i++ {
// 		reqs = append(reqs,
// 			NewRequest(randID(), ele.Index(i).Interface(), method))
// 	}
// 	bs := encodeMultiRequest(reqs)
// 	defer c.conn.Close()
// 	respStr := c.send(string(bs))
// 	resps := parseMultiResponse(respStr)

// 	// replys must be pointer, so it can be set this
// 	eleReply := reflect.ValueOf(replys)
// 	typReply := reflect.TypeOf(replys)
// 	if typReply.Kind() != reflect.Ptr {
// 		return errMultiReplyTypePtr
// 	}

// 	eleReply = eleReply.Elem()
// 	// fill response.Result into replys
// 	for idx, resp := range resps {
// 		convert(resp.Result, eleReply.Index(idx).Interface())
// 	}
// 	return nil
// }
