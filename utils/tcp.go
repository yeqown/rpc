package utils

import (
	"bufio"
	"log"
	"net"
)

// WriteServerTCP ...
func WriteServerTCP(conn net.Conn, byts []byte) error {
	if conn == nil {
		panic("conn is nil")
	}
	// spliter symbol '\n'
	byts = append(byts, byte('\n'))
	if _, err := conn.Write(byts); err != nil {
		log.Printf("conn.Write err: %v\n", err)
		return err
	}

	return nil
}

// WriteClientTCP ...
func WriteClientTCP(conn net.Conn, data []byte) ([]byte, error) {
	// add spliter symbol '\n'
	data = append(data, '\n')
	// send data to server
	if _, err := conn.Write(data); err != nil {
		log.Printf("conn.Write err: %v\n", err)
		return nil, nil
	}

	// get response form server
	message, err := bufio.NewReader(conn).ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return message, nil
}
