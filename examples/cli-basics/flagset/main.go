package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no command provided")
		os.Exit(2)
	}
	cmd := os.Args[1]
	switch cmd {
	case "greet":
		greetCmd := flag.NewFlagSet("greet", flag.ExitOnError)
		msg := greetCmd.String("msg", "CLI Basics", "greeting message")
		err := greetCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("hello and welcome: %s\n", *msg)
	default:
		fmt.Printf("unknown command: %s\n", cmd)
		os.Exit(2)
	}
}
