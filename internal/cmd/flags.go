package cmd

import (
	"errors"
	"flag"
)

// Flags is a helpful struct that contains the global flags for the equivalence
// test binary.
type Flags struct {
	// The relative or absolute path to the directory that contains the golden
	// files.
	GoldenFilesDirectory string

	// The relative or absolute path to the directory that contains the test
	// files and specifications.
	TestingFilesDirectory string

	// The relative or absolute path to the target Terraform binary.
	TerraformBinaryPath string
}

func ParseFlags(command string, args []string) (*Flags, error) {
	fs := flag.NewFlagSet(command, flag.ContinueOnError)

	flags := Flags{}
	fs.StringVar(&flags.GoldenFilesDirectory, "goldens", "", "Absolute or relative path to the directory containing the golden files.")
	fs.StringVar(&flags.TestingFilesDirectory, "tests", "", "Absolute or relative path to the directory containing the tests and specifications.")
	fs.StringVar(&flags.TerraformBinaryPath, "binary", "terraform", "Absolute or relative path to the target Terraform binary.")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if len(flags.GoldenFilesDirectory) == 0 {
		return nil, errors.New("--goldens flag is required")
	}

	if len(flags.TestingFilesDirectory) == 0 {
		return nil, errors.New("--tests flag is required")
	}

	return &flags, nil
}
