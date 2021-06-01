package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mrityunjaygr8/go-ci/dispatcher"
	"github.com/mrityunjaygr8/go-ci/observer"
)

func main() {
	fmt.Println("yo")
	observerCmd := flag.NewFlagSet("observer", flag.ExitOnError)
	observerDir := observerCmd.String("path", "", "path")
	observerServer := observerCmd.String("dispatcher-server", "127.0.0.1:8888", "dispatcher server")

	dispatcherCmd := flag.NewFlagSet("dispatcher", flag.ExitOnError)
	dispatcherHost := dispatcherCmd.String("host", "127.0.0.1", "host")
	dispatcherPort := dispatcherCmd.String("port", "8888", "port")

	if len(os.Args) < 2 {
		fmt.Println("expected one of \"observer\", \"dispatcher\" or \"runner\" commands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "observer":
		observerCmd.Parse(os.Args[2:])
		fmt.Println("subcommand \"observer\"")
		observer.Observe(*observerDir, *observerServer)
	case "dispatcher":
		dispatcherCmd.Parse(os.Args[2:])
		fmt.Println("subcommand \"dispatcher\"")
		dispatcher.Dispatch(*dispatcherHost, *dispatcherPort)
	default:
		fmt.Println("expected one of \"observer\", \"dispatcher\" or \"runner\" commands")
		os.Exit(1)

	}
}
