package rpc

import (
	"errors"
	"fmt"
	"net"
)

var (
	errMultiReplyTypePtr = errors.New("multi reply should be arrry or slice pointer")
	errEmptyCodec        = errors.New("client has an empty codec")
)

// NewClientWithCodec generate a Client, if codec is nil will
// use default gobCodec
func NewClientWithCodec(addr string, codec Codec) *Client {
	if codec == nil {
		codec = newGobCodec()
	}

	return &Client{
		addr:  addr,
		codec: codec,
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
