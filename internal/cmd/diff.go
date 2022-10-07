package cmd

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"

	"github.com/hashicorp/terraform-equivalence-testing/internal/terraform"
	"github.com/hashicorp/terraform-equivalence-testing/internal/tests"
)

func DiffCommandFactory(ui cli.Ui) cli.CommandFactory {
	return func() (cli.Command, error) {
		return &diffCommand{
			ui: ui,
		}, nil
	}
}

type diffCommand struct {
	ui cli.Ui
}

func (cmd *diffCommand) Help() string {
	return strings.TrimSpace(`
Usage: terraform-equivalence-testing diff --goldens=examples/example_golden_files --tests=examples/example_test_cases [--binary=terraform] [--filters=complex_resource,simple_resource]

Compare and report the diff between a fresh run of the equivalence tests and the golden files.

This command will execute all the test cases within the tests directory, and report any differences between the output and the existing golden files.
`)
}

func (cmd *diffCommand) Run(args []string) int {
	flags, err := ParseFlags("diff", args)
	if err != nil {
		cmd.ui.Error(err.Error())
		return 1
	}

	tf, err := terraform.New(flags.TerraformBinaryPath)
	if err != nil {
		cmd.ui.Error(err.Error())
		return 1
	}
	cmd.ui.Output(fmt.Sprintf("Finding diffs in equivalence tests using Terraform v%s with command `%s`", tf.Version(), flags.TerraformBinaryPath))

	testCases, err := tests.ReadFrom(flags.TestingFilesDirectory, flags.TestFilters...)
	if err != nil {
		cmd.ui.Error(err.Error())
		return 1
	}
	cmd.ui.Output(fmt.Sprintf("Found %d test cases in %s\n", len(testCases), flags.TestingFilesDirectory))

	successfulTests := 0
	testsWithDiffs := 0
	failedTests := 0

	for _, test := range testCases {
		cmd.ui.Output(fmt.Sprintf("[%s]: starting...", test.Name))

		output, err := test.RunWith(tf)
		if err != nil {
			failedTests++
			if tfErr, ok := err.(terraform.Error); ok {
				cmd.ui.Output(fmt.Sprintf("[%s]: %s", test.Name, tfErr))
				continue
			}
			cmd.ui.Output(fmt.Sprintf("[%s]: unknown error (%v)", test.Name, err))
			continue
		}

		cmd.ui.Output(fmt.Sprintf("[%s]: computing diffs...", test.Name))

		files, err := output.ComputeDiff(flags.GoldenFilesDirectory)
		if err != nil {
			failedTests++
			cmd.ui.Output(fmt.Sprintf("[%s]: unknown error (%v)", test.Name, err))
			continue
		}

		newFileCount := 0
		noChangeCount := 0
		changeCount := 0

		for file, diff := range files {
			switch diff {
			case tests.NewFile:
				newFileCount++
				cmd.ui.Output(fmt.Sprintf("[%s]: %s was a new file", test.Name, file))
			case tests.NoChange:
				noChangeCount++
				cmd.ui.Output(fmt.Sprintf("[%s]: %s had no diffs", test.Name, file))
			default:
				changeCount++
				cmd.ui.Output(fmt.Sprintf("[%s]: %s had diffs (-want +got):\n%s", test.Name, file, diff))
			}
		}

		successfulTests++
		if newFileCount+changeCount > 0 {
			testsWithDiffs++
		}

		cmd.ui.Output(fmt.Sprintf("[%s]: complete\n", test.Name))
	}

	cmd.ui.Output(fmt.Sprintf("Equivalence testing complete."))
	cmd.ui.Output(fmt.Sprintf("\tAttempted %d test(s).", len(testCases)))

	if successfulTests > 0 {
		cmd.ui.Output(fmt.Sprintf("\t%d test(s) were successful.", successfulTests))
	}

	if testsWithDiffs > 0 {
		cmd.ui.Output(fmt.Sprintf("\t%d test(s) had diffs.", testsWithDiffs))
	}

	if failedTests > 0 {
		cmd.ui.Output(fmt.Sprintf("\t%d test(s) failed.", failedTests))
		return 1
	}

	return 0
}

func (cmd *diffCommand) Synopsis() string {
	return "Compare and report the diff between a fresh run of the equivalence tests and the golden files."
}
