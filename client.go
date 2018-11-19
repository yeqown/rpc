// Package rpc support json-rpc
package rpc

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"reflect"
)

var (
	errMultiReplyTypePtr = errors.New("multi reply should be arrry or slice pointer")
)

// NewClient ... maybe noneed this function
func NewClient() *Client {
	return &Client{}
}

// Client is json-rpc Client
// includes conn field
type Client struct {
	conn net.Conn
}

// Call Client Call Remote method
func (c *Client) Call(id, method string, args, reply interface{}) error {
	defer c.conn.Close()

	if id == "" {
		id = randID()
	}

	req := NewRequest(id, args, method)
	bs := encodeRequest(req)
	respStr := c.send(string(bs))
	resp := parseResponse(respStr)
	convert(resp.Result, reply)
	return nil
}

// CallMulti ...
// TODO: handle with multi params and multi response? and how
func (c *Client) CallMulti(method string, params, replys interface{}) error {
	ele := reflect.ValueOf(params)
	typ := reflect.TypeOf(params)

	if typ.Kind() == reflect.Ptr {
		ele = ele.Elem()
		typ = ele.Type()
	}

	if typ.Kind() != reflect.Slice && typ.Kind() != reflect.Array {
		err := fmt.Errorf("Error: params type %s is not array type or slice", typ.Kind().String())
		return err
	}

	reqs := make([]*Request, 0)
	for i := 0; i < ele.Len(); i++ {
		reqs = append(reqs,
			NewRequest(randID(), ele.Index(i).Interface(), method))
	}
	bs := encodeMultiRequest(reqs)
	defer c.conn.Close()
	respStr := c.send(string(bs))
	resps := parseMultiResponse(respStr)

	// replys must be pointer, so it can be set this
	eleReply := reflect.ValueOf(replys)
	typReply := reflect.TypeOf(replys)
	if typReply.Kind() != reflect.Ptr {
		return errMultiReplyTypePtr
	}

	eleReply = eleReply.Elem()
	// fill response.Result into replys
	for idx, resp := range resps {
		convert(resp.Result, eleReply.Index(idx).Interface())
	}
	return nil
}

// DialTCP to Dial serevr over TCP
func (c *Client) DialTCP(addr string) {
	var err error
	c.conn, err = net.Dial("tcp", addr)
	if err != nil {
		panic(fmt.Errorf("net.Dial tcp get err: %v", err))
	}
}

func (c *Client) send(s string) string {
	fmt.Fprintf(c.conn, s+"\n")
	message, _ := bufio.NewReader(c.conn).ReadString('\n')
	return message
}
