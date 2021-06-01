package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mrityunjaygr8/go-ci/dispatcher"
	"github.com/mrityunjaygr8/go-ci/observer"
	"github.com/mrityunjaygr8/go-ci/runner"
)

func main() {
	fmt.Println("yo")
	observerCmd := flag.NewFlagSet("observer", flag.ExitOnError)
	observerDir := observerCmd.String("path", "", "path")
	observerServer := observerCmd.String("dispatcher-server", "127.0.0.1:8888", "dispatcher server")

	dispatcherCmd := flag.NewFlagSet("dispatcher", flag.ExitOnError)
	dispatcherHost := dispatcherCmd.String("host", "127.0.0.1", "host")
	dispatcherPort := dispatcherCmd.String("port", "8888", "port")

	runnerCmd := flag.NewFlagSet("runner", flag.ExitOnError)
	runnerHost := runnerCmd.String("host", "127.0.0.1", "host")
	runnerPort := runnerCmd.String("port", "", "port")
	runnerDispatcher := runnerCmd.String("dispatcher-server", "127.0.0.1:8888", "dispatcher server")
	runnerRepo := runnerCmd.String("path", "", "path")

	if len(os.Args) < 2 {
		fmt.Println("expected one of \"observer\", \"dispatcher\" or \"runner\" commands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "observer":
		observerCmd.Parse(os.Args[2:])
		fmt.Println("subcommand \"observer\"")
		if *observerDir == "" {
			fmt.Println("Observe directory cannot be blank")
			os.Exit(1)
		}
		observer.Observe(*observerDir, *observerServer)
	case "dispatcher":
		dispatcherCmd.Parse(os.Args[2:])
		fmt.Println("subcommand \"dispatcher\"")
		dispatcher.Dispatch(*dispatcherHost, *dispatcherPort)
	case "runner":
		runnerCmd.Parse(os.Args[2:])
		fmt.Println("Subcommand \"runner\"")
		if *runnerRepo == "" {
			fmt.Println("Runner directory cannot be blank")
			os.Exit(1)
		}
		runner.Run(*runnerHost, *runnerPort, *runnerDispatcher, *runnerRepo)
	default:
		fmt.Println("expected one of \"observer\", \"dispatcher\" or \"runner\" commands")
		os.Exit(1)

	}
}
