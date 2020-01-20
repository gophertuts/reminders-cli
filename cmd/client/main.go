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
	s := client.NewSwitch(*backendURIFlag)

	if *helpFlag || len(os.Args) == 1 {
		s.Help()
		return
	}

	err := s.Switch()
	if err != nil {
		fmt.Printf("cmd switch error: %v\n", err)
		os.Exit(2)
	}
}
