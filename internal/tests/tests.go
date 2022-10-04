package tests

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/hashicorp/terraform-equivalence-testing/internal/files"
	"github.com/hashicorp/terraform-equivalence-testing/internal/terraform"
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

// ReadFrom accepts a directory and returns the set of test cases specified
// within this directory.
func ReadFrom(directory string) ([]Test, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var tests []Test
	for _, file := range files {
		if file.IsDir() {
			data, err := ioutil.ReadFile(path.Join(directory, file.Name(), "spec.json"))
			if err != nil {
				return nil, err
			}

			var specification TestSpecification
			if err := json.Unmarshal(data, &specification); err != nil {
				return nil, err
			}

			tests = append(tests, Test{
				Name:          file.Name(),
				Specification: specification,
				Directory:     directory,
			})
		}
	}
	return tests, nil
}

// RunWith executes the specified test using the Terraform binary specified by
// the terraform.Terraform argument.
//
// This function will return a TestOutput struct, which contains the file names
// of the outputs that we want to compare. These files are already read in and
// parsed in JSON objects.
func (test Test) RunWith(tf terraform.Terraform) (TestOutput, error) {
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
