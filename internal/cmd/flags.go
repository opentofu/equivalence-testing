// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cmd

import (
	"errors"
	"flag"
	"os"
	"path"
	"path/filepath"
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

	// If empty, then all tests will be executed. If not empty, only tests
	// included in this flag will be executed.
	TestFilters StringList
}

func ParseFlags(command string, args []string) (*Flags, error) {
	fs := flag.NewFlagSet(command, flag.ContinueOnError)

	flags := Flags{}

	fs.StringVar(&flags.GoldenFilesDirectory, "goldens", "", "Absolute or relative path to the directory containing the golden files.")
	fs.StringVar(&flags.TestingFilesDirectory, "tests", "", "Absolute or relative path to the directory containing the tests and specifications.")
	fs.StringVar(&flags.TerraformBinaryPath, "binary", "terraform", "Absolute or relative path to the target Terraform binary.")

	fs.Var(&flags.TestFilters, "filters", "If specified, only test cases included in this list will be executed.")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if len(flags.GoldenFilesDirectory) == 0 {
		return nil, errors.New("--goldens flag is required")
	}

	if len(flags.TestingFilesDirectory) == 0 {
		return nil, errors.New("--tests flag is required")
	}

	// Last thing, let's change the TerraformBinaryPath into an absolute path as
	// we are messing around with the working directory later. One exception is
	// if the caller has asked to just execute the default Terraform system
	// command/binary.
	if !filepath.IsAbs(flags.TerraformBinaryPath) && flags.TerraformBinaryPath != "terraform" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		flags.TerraformBinaryPath = path.Join(wd, flags.TerraformBinaryPath)
	}

	return &flags, nil
}
