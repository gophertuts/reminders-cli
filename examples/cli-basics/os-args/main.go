package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("no command provided")
		os.Exit(2)
	}
	cmd := os.Args[1]
	switch cmd {
	case "greet":
		msg := "CLI Basics"
		if len(os.Args) > 2 {
			if os.Args[2] == "--help" {
				fmt.Println("Usage of greet:\n -msg string\n\tgreeting message (default 'CLI Basics')")
				return
			}
			msgFlag := getFlag(os.Args[2], "--msg")
			if msgFlag != "" {
				msg = msgFlag
			}
		}
		fmt.Printf("hello and welcome: %s\n", msg)
	case "help":
		fmt.Println("here's some help, take it")
	default:
		fmt.Printf("unknown command: %s\n", cmd)
		os.Exit(2)
	}
}

func getFlag(s string, name string) string {
	f := strings.Split(s, "=")
	if len(f) != 2 || f[0] != name {
		return ""
	}
	return f[1]
}
