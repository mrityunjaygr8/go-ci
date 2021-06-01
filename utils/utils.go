package utils

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const OK = "ok"
const INVALID_COMMAND = "Invalid Command"
const COMMIT_FILE = "./.commit_id"
const PING = "ping"
const PONG = "pong"
const BUF_SIZE = 2048
const TEST_RESULTS_DIR = "test_results"
const BUSY = "BUSY"
const PORT_RANGE_START = 8900
const PORT_RANGE_END = 9000

func Communicate(host HP, msg string) (string, error) {
	resp := make([]byte, BUF_SIZE)
	conn, err := net.Dial("tcp", host.to_address())
	if err != nil {
		fmt.Println("Error Connecting: ", err.Error())
		return "", err
	}
	defer conn.Close()

	fmt.Fprint(conn, msg+"\a")

	n, err := bufio.NewReader(conn).Read(resp)
	if err != nil {
		fmt.Println("Could not communicate with server")
		return "", err
	}
	return string(resp[:n]), nil
}

func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

type HP struct {
	Host string
	Port string
}

func (h HP) to_address() string {
	return h.Host + ":" + h.Port
}

func HPFromString(hp string) HP {
	res := strings.Split(hp, ":")
	return HP{Host: res[0], Port: res[1]}
}
