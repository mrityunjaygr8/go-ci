package main
import (
	"flag"
	"fmt"
	"os"
	"github.com/mrityunjaygr8/go-ci/observer"
)
func main() {
	fmt.Println("yo")
	observerCmd := flag.NewFlagSet("observer", flag.ExitOnError)
	observerDir := observerCmd.String("path", "", "path")
	observerServer := observerCmd.String("dispatcher-server", "127.0.0.1:8888", "dispatcher server")
	if len(os.Args) < 2 {
		fmt.Println("expected the \"observer\" command")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "observer":
		observerCmd.Parse(os.Args[2:])
		fmt.Println("subcommand \"observer\"")
		observer.Observe(*observerDir, *observerServer)
	default:
		fmt.Println("expected the \"observer\" command")
		os.Exit(1)

	}
}
