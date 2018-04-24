package rpc

import (
	"bufio"
	"fmt"
	"net"
)

type Client struct {
	conn net.Conn
}

func (c *Client) Call() {

}

func (c *Client) Go() {

}

func (c *Client) DialTCP(addr string) {
	var err error
	c.conn, err = net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connect Ok")
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
