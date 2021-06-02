package dispatcher

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/mrityunjaygr8/go-ci/utils"
)

type Dispatcher struct {
	runners           *[]utils.HP
	dispatched_commit map[string]utils.HP
	pending           *[]string
	server            utils.CTS
	dead              bool
}

func (d *Dispatcher) handleDispatcher(conn net.Conn, message string) {
	re := regexp.MustCompile(`(?P<command>\w+)(?P<followup>:.+)*`)
	if !re.Match([]byte(message)) {
		conn.Write([]byte(utils.INVALID_COMMAND))
	}
	sub := re.FindStringSubmatch(message)
	command := sub[1]
	fmt.Println("Command Recieved: ", command)
	switch command {
	case "status":
		fmt.Println("in status")
		conn.Write([]byte(utils.OK))

	case "register":
		followup := sub[2]
		register_re := regexp.MustCompile(`:([\w.]*)`)
		register_sub := register_re.FindAllStringSubmatch(followup, -1)
		host, port := register_sub[0][1], register_sub[1][1]
		*d.runners = append(*d.runners, utils.HP{Host: host, Port: port})
		fmt.Println(d.runners)
		conn.Write([]byte(utils.OK))

	case "dispatch":
		commit_id := sub[2][1:]
		fmt.Println(d.pending)
		fmt.Println("Going to dispatch")
		if len(*d.runners) == 0 {
			conn.Write([]byte("No runners are registered"))
			*d.pending = append(*d.pending, commit_id)
		} else {
			conn.Write([]byte(utils.OK))
			d.dispatchTest(commit_id)
		}
	case "results":
		fmt.Println("we are in results")
		results := strings.Split(sub[2][1:], ":")
		commit_id := results[0]

		prefix := len(command) + len(commit_id) + len(results[1]) + 3
		test_result := message[prefix:]

		delete(d.dispatched_commit, commit_id)
		if _, err := os.Stat(utils.TEST_RESULTS_DIR); os.IsNotExist(err) {
			os.Mkdir(utils.TEST_RESULTS_DIR, 0700)
		}

		err := ioutil.WriteFile(utils.TEST_RESULTS_DIR+"/"+commit_id, []byte(test_result), 0644)
		if err != nil {
			fmt.Println("An error has occurred: ", err)
		}
		conn.Write([]byte(utils.OK))

	default:
		conn.Write([]byte(utils.INVALID_COMMAND))
	}
}

func (d *Dispatcher) manageCommits(runner utils.HP, idx int) {
	for commit, assigned_runner := range d.dispatched_commit {
		if assigned_runner == runner {
			delete(d.dispatched_commit, commit)
			*d.pending = append(*d.pending, commit)
			break
		}
	}
	d.deleteRunners(idx)
}

func (d *Dispatcher) deleteRunners(idx int) {
	runners := *d.runners
	runners[idx] = runners[len(runners)-1]
	*d.runners = runners[:len(runners)-1]
}

func (d *Dispatcher) runnerChecker() {
	for !d.dead {
		fmt.Println("runner checker")
		time.Sleep(1 * time.Second)
		for idx, runner := range *d.runners {
			runner_resp, err := utils.Communicate(runner, utils.PING)
			if err != nil {
				fmt.Println("An error has occurred: ", err)
			}
			if runner_resp != utils.PONG {
				fmt.Println("removing", runner.Host, runner.Port)
				d.manageCommits(runner, idx)
			}
		}
	}
}

func (d *Dispatcher) dispatchTest(commit_id string) {
	fmt.Println("trying to dispatch to runners")
	for {
		for _, runner := range *d.runners {
			runner_resp, err := utils.Communicate(runner, "runtest:"+commit_id)
			if err != nil {
				fmt.Println("An error has occurred: ", err)
			}

			if runner_resp == utils.OK {
				fmt.Println("Adding id " + commit_id)
				d.dispatched_commit[commit_id] = runner

				idx, found := utils.Find(*d.pending, commit_id)

				if found {
					pending := *d.pending
					pending[idx] = pending[len(pending)-1]
					*d.pending = pending[:len(pending)-1]
				}

				return
			} else if runner_resp == utils.BUSY {
				*d.pending = append(*d.pending, commit_id)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func (d *Dispatcher) redistribute() {
	fmt.Println("Redis")
	for !d.dead {
		for _, commit := range *d.pending {
			fmt.Println("Running Redistribute")
			d.dispatchTest(commit)
			time.Sleep(5 * time.Second)
		}
	}
}

func newDispatch(host, port string) Dispatcher {
	dispatcher := Dispatcher{}
	dispatcher.dispatched_commit = make(map[string]utils.HP)
	dispatcher.runners = &[]utils.HP{}
	dispatcher.pending = &[]string{}
	dispatcher.dead = false
	server := utils.NewCTS(host, port, dispatcher.handleDispatcher)
	dispatcher.server = server
	return dispatcher
}

func Dispatch(host, port string) {
	dispatcher := newDispatch(host, port)
	go dispatcher.server.Start()
	go dispatcher.runnerChecker()
	go dispatcher.redistribute()

	for {
		if dispatcher.dead {
			dispatcher.server.Stop()
		}
	}
}
