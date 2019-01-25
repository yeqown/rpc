package rpc

import (
	"errors"
	"fmt"
	"net"

	"github.com/yeqown/rpc/utils"
)

var (
	errMultiReplyTypePtr = errors.New("multi reply should be arrry or slice pointer")
	errEmptyCodec        = errors.New("client has an empty codec")
)

// NewClientWithCodec generate a Client
// prototype rpc.NewClientWithCodec(codec Codec, tcpAddr string, httpAddr string)
// if codec is nil will use default gobCodec, tcpAddr or httpAddr is empty only when
// you are sure about it will never be used, otherwise it panic while using some functions.
func NewClientWithCodec(codec Codec, tcpAddr, httpAddr string) *Client {
	if codec == nil {
		codec = newGobCodec()
	}

	return &Client{
		tcpAddr:  tcpAddr,
		httpAddr: httpAddr,
		codec:    codec,
	}
}

// Client ....
type Client struct {
	// rpc server addr over tcp
	tcpAddr string

	// rpc server addr over http
	httpAddr string

	// codec to manage about the request and response encoding and decoding
	codec Codec

	// connection to the tcp server
	tcpConn net.Conn
}

// Call Client Call Remote method
// TODO: timeout cancel
func (c *Client) Call(method string, args, reply interface{}) error {
	if c.codec == nil {
		return errEmptyCodec
	}

	// connect to server
	if c.tcpConn == nil {
		conn, err := net.Dial("tcp", c.tcpAddr)
		if err != nil {
			return fmt.Errorf("net.Dial tcp get err: %v", err)
		}
		c.tcpConn = conn
	}

	// core ....
	req := c.codec.Request(method, args)
	respDataByts, err := utils.WriteClientTCP(c.tcpConn, encodeRequest(c.codec, req))
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

// CallHTTP generate a http request and send to the server
func (c *Client) CallHTTP(method string, args, reply interface{}) error {
	// assemble the request and encode the request to []byte
	rpcReq := c.codec.Request(method, args)
	data, err := c.codec.Encode(rpcReq)
	if err != nil {
		return fmt.Errorf("c.codec.Encode(rpcReq) got err: %v", err)
	}

	debugF("send request [addr: %s] [data: %s]", c.httpAddr, data)
	// request the server
	byts, err := utils.RequestHTTP(c.httpAddr, data)
	if err != nil {
		return fmt.Errorf("c.codec.Encode(rpcReq) got err: %v", err)
	}
	debugF("got response %s\n", byts)

	var rpcResp Response
	if rpcResp, err = c.codec.ParseResponse(byts); err != nil {
		return err
	}
	if err := rpcResp.Error(); err != nil {
		return fmt.Errorf("rpcResp.Error(): %v", err)
	}
	if err := c.codec.Decode(rpcResp.Reply(), reply); err != nil {
		return fmt.Errorf("c.codec.Decode(rpcResp.Reply() got err: %v", err)
	}

	return nil
}

// Close the client connectio to the server
func (c *Client) Close() {
	if c.tcpConn == nil {
		return
	}
	c.tcpConn.Close()
}

// encodeRequest ...
func encodeRequest(codec Codec, req Request) []byte {
	byts, err := codec.Encode(req)
	if err != nil {
		panic(err)
	}
	return byts
}
