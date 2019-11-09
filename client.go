package rpc

import (
	"bufio"
	"errors"
	"fmt"
	"net"

	"github.com/yeqown/rpc/proto"
	// "github.com/yeqown/rpc/utils"
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

// Call .
func (c *Client) Call(method string, args, reply interface{}) error {
	req := c.codec.NewRequest(method, args)
	resps := make([]Response, 0)
	if err := c.calltcp([]Request{req}, &resps); err != nil {
		DebugF("could not calltcp err=%v", err)
		return err
	}

	resp := resps[0]
	// DebugF("len(resps)=%d, stdResponse=%v", len(resps), resps[0].(*stdResponse))
	DebugF("resp.Reply()=%s", resp.Reply())
	// TODO: handle reponse error
	// TODO: handle resp maybe nil
	if err := c.codec.ReadResponseBody(resp.Reply(), reply); err != nil {
		DebugF("could not ReadReponseBody err=%v", err)
		return err
	}

	return nil
}

// Call server over tcp
// TODO: timeout cancel
func (c *Client) calltcp(reqs []Request, resps *[]Response) (err error) {
	if err = c.valid(); err != nil {
		return err
	}

	var (
		wr    = bufio.NewWriter(c.tcpConn)
		rr    = bufio.NewReader(c.tcpConn)
		psend = proto.New()
		precv = proto.New()
	)

	// req := c.codec.NewRequest(method, args)
	if psend.Body, err = c.codec.EncodeRequests(&reqs); err != nil {
		DebugF("could not EncodeRequests, err=%v", err)
		return err
	}

	if err := psend.WriteTCP(wr); err != nil {
		DebugF("could not WriteTCP, err=%v", err)
		return err
	}
	wr.Flush()

	// recv from TCP response
	if err := precv.ReadTCP(rr); err != nil {
		DebugF("could not ReadTCP, err=%v", err)
		return err
	}

	DebugF("recv response body: %s", precv.Body)
	// var resp Response
	*resps, err = c.codec.ReadResponse(precv.Body)
	if err != nil {
		DebugF("could not ReadResponses, err=%v", err)
		return err
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

func (c *Client) valid() error {
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
