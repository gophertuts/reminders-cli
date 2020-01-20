package main

import (
	"flag"

	"github.com/gophertuts/reminders-cli/client"
)

var (
	backendURIFlag = flag.String("backend", "http://localhost:8080", "Backend API URI")
	helpFlag       = flag.Bool("help", false, "Display a helpful message")
)

func main() {
	flag.Parse()
	s := client.NewSwitch(*helpFlag, *backendURIFlag)
	s.Switch()
}
