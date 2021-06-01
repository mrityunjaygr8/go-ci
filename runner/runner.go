package runner

import (
	"fmt"
	"net"
	"os"

	"github.com/mrityunjaygr8/go-ci/utils"
)

type runner struct {
	server     utils.CTS
	dispatcher utils.HP
	repo       string
	busy       bool
	dead       bool
}

func (r *runner) handleRunner(net.Conn, string) {}

func newRunner(host, port, repo string, dispatcher utils.HP) runner {
	runner := runner{}
	runner.dispatcher = dispatcher
	runner.busy = false
	runner.dead = false
	runner.repo = repo

	server := setupServer(host, port, runner.handleRunner)
	runner.server = server
	return runner

}

func setupServer(host, port string, handler func(net.Conn, string)) utils.CTS {
	if port == "" {
		port := utils.PORT_RANGE_START
		for port < utils.PORT_RANGE_END {
			cts := utils.NewCTS("127.0.0.1", fmt.Sprint(port), handler)
			err := cts.Start()

			if err != nil {
				port++
				continue
			}
			cts.Stop()
			return cts
		}
		fmt.Println("Can not start the server within the defined port range")
		os.Exit(1)
	} else {
		cts := utils.NewCTS(host, port, handler)
		err := cts.Start()
		if err != nil {
			fmt.Println("Can't listen on ", port)
			os.Exit(1)
		}
		return cts
	}
	return utils.CTS{}
}

func Run(host, port, dispatcher, repo string) {
	fmt.Println(host, port, dispatcher, repo)
	fmt.Println("run")

	dispatcher_hp := utils.HPFromString(dispatcher)
	runner := newRunner(host, port, repo, dispatcher_hp)

}
