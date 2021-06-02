package runner

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/mrityunjaygr8/go-ci/utils"
)

type runner struct {
	server             utils.CTS
	dispatcher         utils.HP
	repo               string
	busy               bool
	dead               bool
	last_communication time.Time
}

func (r *runner) handleRunner(conn net.Conn, message string) {
	command_re := regexp.MustCompile(`(?P<command>\w+)(?P<followup>:.+)*`)
	if !command_re.Match([]byte(message)) {
		conn.Write([]byte(utils.INVALID_COMMAND))
	}
	sub := command_re.FindStringSubmatch(message)
	command := sub[1]

	switch command {
	case utils.PING:
		fmt.Println("pinged")
		r.last_communication = time.Now()
		conn.Write([]byte(utils.PONG))
	case "runtest":
		if r.busy {
			conn.Write([]byte(utils.BUSY))
		} else {
			commit_id := sub[2][1:]
			fmt.Println("commit:", commit_id)
			r.busy = true
			stdout, compError := r.runTest(commit_id)
			r.busy = false
			if !compError {
				fmt.Println("An error has occurred ", &stdout)
				conn.Write([]byte(fmt.Sprintf("An error has occurred, %s", &stdout)))
			} else {
				fmt.Println("Output ", &stdout)
				_, err := utils.Communicate(r.dispatcher, fmt.Sprintf("results:%s:%d:%s", commit_id, stdout.Len()+1, stdout.String()))
				if err != nil {
					fmt.Println("An error has occerred ", err)
				}
			}
		}
	default:
		fmt.Println("Invalid Command Recieved: ", message)
		conn.Write([]byte(utils.INVALID_COMMAND))
	}
}

func (r *runner) runTest(commit_id string) (bytes.Buffer, bool) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("./test_runner_script.sh", r.repo, commit_id)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("An error has occurred ", err)
	}
	fmt.Printf("%s", out)
	goExe, err := exec.LookPath("go")
	if err != nil {
		log.Fatal("failed to find go executable: ", err)
	}
	cmd_test := &exec.Cmd{
		Dir:    r.repo,
		Path:   goExe,
		Args:   []string{goExe, "test"},
		Stdout: &stdout,
		Stderr: &stderr,
	}

	err = cmd_test.Run()
	if err == nil {
		return stdout, true
	}

	exitError, ok := err.(*exec.ExitError)
	if !ok {
		fmt.Printf("command \"%s\" failed with non exit error %s", cmd_test.String(), err)
	}

	switch exc := exitError.ExitCode(); exc {
	case 1:
		return stdout, true
	case 2:
		//  go test returns 2 on a compilation / build error
		stdout.WriteString(fmt.Sprintf("'%s' returned exit code %d: %s",
			cmd_test.String(), exc, err,
		))
		return stdout, false
	default:
		fmt.Printf("error: '%s' failed with exit error %d: %s",
			cmd_test.String(), exc, err,
		)
	}

	return stdout, false
}

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
			fmt.Println(port)
			cts := utils.NewCTS(host, fmt.Sprint(port), handler)
			err := cts.Test()

			if err != nil {
				port++
				continue
			}
			return cts
		}
		fmt.Println("Can not start the server within the defined port range")
		os.Exit(1)
	} else {
		cts := utils.NewCTS(host, port, handler)
		err := cts.Test()
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
	fmt.Println("Registering with dispatcher")
	register_resp, err := utils.Communicate(runner.dispatcher, fmt.Sprintf("register:%s", runner.server.Address()))
	if err != nil {
		fmt.Println("An error has occurred ", err)
		fmt.Println("Cannot register with dispatcher. Exiting")
		os.Exit(1)
	}

	if register_resp != utils.OK {
		fmt.Println("Cannot register with dispatcher. Exiting")
		os.Exit(1)
	}

	go runner.server.Start()

	for {
		if runner.dead {
			runner.server.Stop()
		}
	}

}
