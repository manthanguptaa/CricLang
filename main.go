package main

import (
	"CricLang/repl"
	"fmt"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome %s to CricLang: A fun programming language for cricket enthusiasts!\n", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
