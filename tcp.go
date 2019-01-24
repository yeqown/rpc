package rpc

import (
	"bufio"
	"log"
	"net"
)

// WriteServerTCP ...
func WriteServerTCP(conn net.Conn, codec Codec, resp Response) error {
	if conn == nil {
		panic("conn is nil")
	}

	byts, err := codec.Encode(resp)
	if err != nil {
		log.Printf("writeTCP codec.Encode(resp) got err: %v\n", err)
		return err
	}
	// spliter symbol \n
	byts = append(byts, byte('\n'))
	if _, err := conn.Write(byts); err != nil {
		log.Printf("conn.Write err: %v\n", err)
		return err
	}

	return nil
}

// WriteClientTCP ...
func WriteClientTCP(conn net.Conn, codec Codec, req Request) ([]byte, error) {
	// codec to encode request
	data, err := codec.Encode(req)
	if err != nil {
		return nil, err
	}
	debugF("send data: %v and encoded to be: %v", req, data)

	// add spliter symbol
	data = append(data, '\n')

	// send data to server
	if _, err := conn.Write(data); err != nil {
		debugF("conn.Write err: %v", err)
		return nil, nil
	}

	// get response form server
	message, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return message, nil
}
