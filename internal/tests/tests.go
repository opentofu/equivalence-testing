// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import (
	"bytes"
	"os"
	"path"
	"path/filepath"

	"github.com/komkom/jsonc/jsonc"
	"github.com/opentffoundation/equivalence-testing/internal/binary"
	"github.com/opentffoundation/equivalence-testing/internal/files"
)

// Test defines a single equivalence test within our framework.
//
// Each test has a Name that references the directory that contains our testing
// data. Within this directory there should be a `spec.json` file which is
// read in the TestSpecification object.
//
// The Directory variable references the parent directory of the test, so the
// full path for a given test case is paths.Join(test.Directory, test.Name).
type Test struct {
	Name          string
	Directory     string
	Specification TestSpecification
}

func contains(test string, filters []string) bool {
	for _, filter := range filters {
		if filter == test {
			return true
		}
	}
	return false
}

// ReadFrom accepts a directory and returns the set of test cases specified
// within this directory.
func ReadFrom(directory string, globalRewrites map[string]map[string]string, filters ...string) ([]Test, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var tests []Test
	for _, file := range files {
		if file.IsDir() {
			if len(filters) == 0 || contains(file.Name(), filters) {
				data, err := os.ReadFile(path.Join(directory, file.Name(), "spec.json"))
				if err != nil {
					return nil, err
				}

				var specification TestSpecification

				decoder, err := jsonc.NewDecoder(bytes.NewReader(data))
				if err != nil {
					return nil, err
				}

				if err := decoder.Decode(&specification); err != nil {
					return nil, err
				}

				specification.AddRewrites(globalRewrites)

				tests = append(tests, Test{
					Name:          file.Name(),
					Specification: specification,
					Directory:     directory,
				})
			}
		}
	}
	return tests, nil
}

// RunWith executes the specified test using the binary specified by
// the binary.Binary argument.
//
// This function will return a TestOutput struct, which contains the file names
// of the outputs that we want to compare. These files are already read in and
// parsed in JSON objects.
func (test Test) RunWith(tf binary.Binary) (TestOutput, error) {
	tmp, err := os.MkdirTemp(test.Directory, test.Name)
	if err != nil {
		return TestOutput{}, err
	}
	defer os.RemoveAll(tmp)

	testDirectory := path.Join(test.Directory, test.Name)
	if err = filepath.WalkDir(testDirectory, files.CopyDir(testDirectory, tmp, []string{"spec.json"})); err != nil {
		return TestOutput{}, err
	}

	files, err := tf.ExecuteTest(tmp, test.Specification.IncludeFiles, test.Specification.Commands...)

	if err != nil {
		return TestOutput{}, err
	}

	return TestOutput{
		Test:  test,
		files: files,
	}, nil
}
