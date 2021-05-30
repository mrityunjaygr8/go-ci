package observer

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mrityunjaygr8/go-ci/utils"
)

func Observe(path, server string) {
	server_str := strings.Split(server, ":")
	host := utils.HP{Host: server_str[0], Port: server_str[1]}
	for {
		cmd := exec.Command("./update_repo.sh", path+"_obs")
		_, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}

		if _, err := os.Stat(utils.COMMIT_FILE); err == nil {
			fmt.Println("File exists")
			status_resp, err := utils.Communicate(host, "status")
			if err != nil {
				fmt.Println("An error has occurred: ", err)
			}
			if status_resp == utils.OK {
				commit, err := ioutil.ReadFile(utils.COMMIT_FILE)
				if err != nil {
					log.Fatal(err)
				}
				trimmed_commit := strings.TrimSpace(string(commit))
				dispatch_resp, err := utils.Communicate(host, "dispatch:"+trimmed_commit)

				if err != nil {
					fmt.Println("An error has occurred: ", err)
				}

				if dispatch_resp != utils.OK {
					fmt.Println("could not dispatch the test: ", dispatch_resp)
				}
				fmt.Println("dispatched")
			}
		} else if os.IsNotExist(err) {
			fmt.Println("File does not exits")
		} else {
			log.Fatal(err)
		}

		time.Sleep(5 * time.Second)
	}
}
