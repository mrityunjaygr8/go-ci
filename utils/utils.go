package utils
const OK = "ok"
const COMMIT_FILE = "./.commit_id"
const BUF_SIZE = 2048
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

type HP struct {
	Host string
	Port string
}

func (h HP) to_address() string {
	return h.Host + ":" + h.Port
}
