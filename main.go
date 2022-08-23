package main

import (
	"github.com/jooyyy/larkPartner/watcher"
	"github.com/mitchellh/cli"
	"log"
	"os"
)

func main() {
	c := cli.NewCLI("watcher", "1.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"watcher": watcher.CommandFactory,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
