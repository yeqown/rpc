// Package rpc support json-rpc
package rpc

import (
	"bufio"
	"fmt"
	"net"
)

func NewClient() *Client {
	return &Client{}
}

// Clien is json-rpc Client
// includes conn field
type Client struct {
	conn net.Conn
}

// Client Call Remote method
func (c *Client) Call(id, method string, args, reply interface{}) error {
	defer c.conn.Close()

	if id == "" {
		id = randId()
	}

	req := NewRequest(id, args, method)
	bs := encodeRequest(req)
	resp_s := c.send(string(bs))
	resp := parseResponse(resp_s)
	convert(resp.Result, reply)
	return nil
}

// DialTCP to Dial serevr over TCP
func (c *Client) DialTCP(addr string) {
	var err error
	c.conn, err = net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
}

func (c *Client) send(s string) string {
	fmt.Fprintf(c.conn, s+"\n")
	message, _ := bufio.NewReader(c.conn).ReadString('\n')
	return message
}
