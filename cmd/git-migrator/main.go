package main

import (
	"os"

	"github.com/adamf123git/git-migrator/cmd/git-migrator/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
