/*
 * codec:
 * param encode and decode
 * this is a copy from golang.org
 */

package rpc

import (
	"bufio"
	"encoding/gob"
	"io"
)

type gobServerCodec struct {
	rwc    io.ReadWriteCloser
	dec    *gob.Decoder
	enc    *gob.Encoder
	encBuf *bufio.Writer
	closed bool
}

func (c *gobServerCodec) ReadRequest(body interface{}) error {
	return c.dec.Decode(body)
}

func (c *gobServerCodec) ReadResponse(resp interface{}) error {
	return nil
}

func (c *gobServerCodec) WriteResponse(resp interface{}) error {
	return nil
}

func (c *gobServerCodec) Close() error {
	if c.closed {
		return nil
	}
	c.closed = true
	return c.rwc.Close()
}
