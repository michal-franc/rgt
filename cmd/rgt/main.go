package main

import (
	"github.com/michal-franc/rgt/internal/app/rgt/commands"
	"os"
)

const DefaultCommand = "start"

func main() {
	// if no command specified add in default command
	if len(os.Args) <= 1 {
		os.Args = append(os.Args, DefaultCommand)
	}
	commands.Execute()
}
