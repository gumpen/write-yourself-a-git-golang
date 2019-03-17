package main

import (
	"log"
	"os"

	"github.com/gumpen/write-yourself-a-git-golang"
	"github.com/mitchellh/cli"
)

// cli.Command interfaceをサブコマンドごとに実装する必要がある
// それらをc.Commandsに登録してc.Run()
func main() {
	c := cli.NewCLI("app", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"add": func() (cli.Command, error) {
			return &wyag.AddCommand{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
