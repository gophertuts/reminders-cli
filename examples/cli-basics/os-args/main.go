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
			f := strings.Split(os.Args[2], "=")
			if len(f) == 2 && f[0] == "--msg" {
				msg = f[1]
			}
		}
		fmt.Printf("hello and welcome: %s\n", msg)
	case "help":
		fmt.Println("some help message")
	default:
		fmt.Printf("unknown command: %s\n", cmd)
		os.Exit(2)
	}
}
