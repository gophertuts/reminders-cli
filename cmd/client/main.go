package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gophertuts/reminders-cli/client"
)

var (
	backendURIFlag = flag.String("backend", "http://localhost:8080", "Backend API URI")
	helpFlag       = flag.Bool("help", false, "Display a helpful message")
)

func main() {
	flag.Parse()
	if *helpFlag || len(os.Args) == 1 {
		fmt.Print(client.HelpMsg)
		return
	}
	s := client.NewSwitch(os.Args, *backendURIFlag)
	s.Switch()
}
