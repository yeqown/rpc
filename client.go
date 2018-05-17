// Package rpc support json-rpc
package rpc

import (
	"bufio"
	"fmt"
	"net"
)

// Clien is json-rpc Client
// includes conn field
type Client struct {
	conn net.Conn
}

// TODO: randId
func randId() string {
	return "rand_id"
}

// Client Call Remote method
func (c *Client) Call(id, method string, args, reply interface{}) error {
	req := NewRequest(id, args, method)
	bs := encodeRequest(req)
	resp_s := c.Send(string(bs))
	// print(resp_s)
	resp := parseResponse(resp_s)
	convert(resp.Result, reply)
	return nil
}

func (c *Client) DialTCP(addr string) {
	var err error
	c.conn, err = net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	// fmt.Println("Connect Ok")
}

func (c *Client) Send(s string) string {
	// send
	fmt.Fprintf(c.conn, s+"\n")
	// receive
	message, _ := bufio.NewReader(c.conn).ReadString('\n')
	return message
}

func (c *Client) CloseTCP() {
	c.conn.Close()
}

func NewClient() *Client {
	return &Client{}
}
