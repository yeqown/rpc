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
	errNotSupportMulti   = errors.New("current codec not support multi request")
)

// NewClientWithCodec generate a Client
// prototype rpc.NewClientWithCodec(codec Codec, tcpAddr string, httpAddr string)
// if codec is nil will use default gobCodec, tcpAddr or httpAddr is empty only when
// you are sure about it will never be used, otherwise it panic while using some functions.
func NewClientWithCodec(codec Codec, tcpAddr, httpAddr string) *Client {
	if codec == nil {
		codec = NewGobCodec()
	}

	return &Client{
		tcpAddr:  tcpAddr,
		httpAddr: httpAddr,
		codec:    codec,
	}
}

// Client as a data struct to connect to server, send and recv data
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

// CallOverTCP call server over tcp
// TODO: timeout cancel
func (c *Client) CallOverTCP(method string, args, reply interface{}) error {
	if err := c.validSelf(); err != nil {
		return err
	}

	// core ....
	req := c.codec.Request(method, args)
	tcpRespByts, err := utils.WriteClientTCP(c.tcpConn, encodeRequest(c.codec, req))
	if err != nil {
		return fmt.Errorf("c.codec.Request(c.conn, req) got err: %v", err)
	}

	var resp Response
	if resp, err = c.codec.ParseResponse(tcpRespByts); err != nil {
		return err
	}
	if err := resp.Error(); err != nil {
		return fmt.Errorf("resp.Error(): %v", err)
	}
	// core ...

	if err := c.codec.Decode(resp.Reply(c.codec), reply); err != nil {
		return fmt.Errorf("c.codec.Decode(resp.Reply() got err: %v", err)
	}
	return nil
}

// CallOverTCPMulti ...
func (c *Client) CallOverTCPMulti(rpcReqs []*RequestConfig) error {
	if err := c.validSelf(); err != nil {
		return err
	}
	if !c.codec.MultiSupported() {
		return errNotSupportMulti
	}

	// assemble multi Request into a Request
	reqMulti := c.codec.RequestMulti(rpcReqs)
	data := encodeRequest(c.codec, reqMulti)
	debugF("CallOverTCPMulti encodeRequest(c.codec, reqMulti): %s", data)
	tcpRespByts, err := utils.WriteClientTCP(c.tcpConn, data)
	if err != nil {
		return fmt.Errorf("c.codec.Request(c.conn, req) got err: %v", err)
	}

	var resp Response
	if resp, err = c.codec.ParseResponse(tcpRespByts); err != nil {
		return err
	}
	debugF("CallOverTCPMulti: origin %s, decoded: %v", tcpRespByts, resp)

	counter := 0
	for resp.HasNext() {
		next := resp.Next().(Response)
		debugF("client get response: %v, %s", next, next.Reply(c.codec))
		if err := c.codec.Decode(next.Reply(c.codec), rpcReqs[counter].Reply); err != nil {
			return fmt.Errorf("c.codec.Decode(resp.Reply() got err: %v", err)
		}
		counter++
	}
	debugF("CallOverTCPMulti result: %v", rpcReqs)
	return nil
}

// CallOverHTTP generate a http request and send to the server
func (c *Client) CallOverHTTP(method string, args, reply interface{}) error {
	if err := c.validSelf(); err != nil {
		return err
	}

	// assemble the request and encode the request to []byte
	rpcReq := c.codec.Request(method, args)
	reqEncoded, err := c.codec.Encode(rpcReq)
	if err != nil {
		return fmt.Errorf("c.codec.Encode(rpcReq) got err: %v", err)
	}

	debugF("send request [addr: %s] [reqEncoded: %s]", c.httpAddr, reqEncoded)
	// request the server
	httpRespByts, err := utils.RequestHTTP(c.httpAddr, reqEncoded)
	if err != nil {
		return fmt.Errorf("c.codec.Encode(rpcReq) got err: %v", err)
	}
	debugF("got response %s\n", httpRespByts)

	var rpcResp Response
	if rpcResp, err = c.codec.ParseResponse(httpRespByts); err != nil {
		return err
	}
	if err := rpcResp.Error(); err != nil {
		return fmt.Errorf("rpcResp.Error(): %v", err)
	}
	debugF("client decode to reply: origin %s, decoded: %v", rpcResp.Reply(c.codec), reply)
	if err := c.codec.Decode(rpcResp.Reply(c.codec), reply); err != nil {
		return fmt.Errorf("c.codec.Decode(rpcResp.Reply() got err: %v", err)
	}
	return nil
}

// CallOverHTTPMulti ...
func (c *Client) CallOverHTTPMulti(rpcReqs []*RequestConfig) error {
	if err := c.validSelf(); err != nil {
		return err
	}

	if !c.codec.MultiSupported() {
		return errNotSupportMulti
	}

	// assemble multi Request into a Request
	reqMulti := c.codec.RequestMulti(rpcReqs)
	reqEncoded, err := c.codec.Encode(reqMulti)
	if err != nil {
		return fmt.Errorf("c.codec.Encode(reqMulti) got err: %v", err)
	}

	debugF("send request [addr: %s] [reqEncoded: %s]", c.httpAddr, reqEncoded)
	// request the server
	httpRespByts, err := utils.RequestHTTP(c.httpAddr, reqEncoded)
	if err != nil {
		return fmt.Errorf("c.codec.Encode(rpcReq) got err: %v", err)
	}
	debugF("got response %s\n", httpRespByts)

	var rpcResp Response
	if rpcResp, err = c.codec.ParseResponse(httpRespByts); err != nil {
		return err
	}

	counter := 0
	for rpcResp.HasNext() {
		next := rpcResp.Next().(Response)
		debugF("client decode to reply: origin %s, decoded: %v", rpcResp.Reply(c.codec), rpcReqs[counter].Reply)
		if err := c.codec.Decode(next.Reply(c.codec), rpcReqs[counter].Reply); err != nil {
			return fmt.Errorf("c.codec.Decode(v.Reply(), v) got err: %v", err)
		}
		counter++
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

func (c *Client) validSelf() error {
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
	return nil
}

// encodeRequest ...
func encodeRequest(codec Codec, req Request) []byte {
	byts, err := codec.Encode(req)
	if err != nil {
		panic(err)
	}
	return byts
}
