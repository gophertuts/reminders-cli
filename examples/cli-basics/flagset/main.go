package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// Commands & Flags
const (
	GreetCmd = "greet"
	MsgFlag  = "msg"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("no command provided")
	}
	cmd := os.Args[1]
	switch cmd {
	case GreetCmd:
		greetCmd := flag.NewFlagSet(GreetCmd, flag.ExitOnError)
		msgFlag := greetCmd.String(MsgFlag, "CLI Basics", "greeting message")
		err := greetCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("hello and welcome: %s\n", *msgFlag)
	default:
		log.Fatalf("unknown command: %s\n", cmd)
	}
}
