package utils

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type CTS struct {
	Handler func(net.Conn, string)
	address HP
}

func NewCTS(host, port string, handler func(net.Conn, string)) CTS {
	c := CTS{Handler: handler, address: HP{Host: host, Port: port}}
	return c
}

func (c *CTS) Start() {
	ln, err := net.Listen("tcp", c.address.Host+":"+c.address.Port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go c.handle(conn, c.Handler)
	}
}

func (c *CTS) Stop() {
	fmt.Println("Exiting")

}

func (c *CTS) handle(conn net.Conn, handler func(net.Conn, string)) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	scanner.Split(customSplitFunc)
	for scanner.Scan() {
		message := scanner.Text()
		handler(conn, message)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error: ", err)
	}
}

func customSplitFunc(data []byte, atEOF bool) (advance int, token []byte, err error) {

	// Return nothing if at end of file and no data passed
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	// Find the index of the input of a newline followed by a
	// pound sign.
	if i := strings.Index(string(data), "\a"); i >= 0 {
		return i + 1, data[0:i], nil
	}

	// If at end of file with data return the data
	if atEOF {
		return len(data), data, nil
	}

	return
}
