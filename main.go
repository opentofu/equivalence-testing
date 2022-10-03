package main

import (
	"fmt"
	"os"

	"github.com/mitchellh/cli"

	"github.com/hashicorp/terraform-equivalence-testing/internal/cmd"
)

var (
	version = "dev"
)

func main() {
	ui := cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	command := cli.NewCLI("terraform-equivalence-testing", version)

	command.Args = os.Args[1:]
	command.Commands = map[string]cli.CommandFactory{
		"diff":   cmd.DiffCommandFactory(&ui),
		"update": cmd.UpdateCommandFactory(&ui),
	}
	command.HelpFunc = cli.BasicHelpFunc("terraform-equivalence-testing")
	command.HelpWriter = os.Stdout
	command.ErrorWriter = os.Stderr

	status, err := command.Run()
	if err != nil {
		ui.Error(fmt.Sprintf("failed to execute equivalence tests: %v", err))
	}
	os.Exit(status)
}
