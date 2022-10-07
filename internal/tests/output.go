package tests

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"path/filepath"

	"github.com/google/go-cmp/cmp"

	"github.com/hashicorp/terraform-equivalence-testing/internal/files"
	strip "github.com/hashicorp/terraform-equivalence-testing/internal/json"
)

const (
	NewFile  string = "(new file)"
	NoChange string = "(no change)"
)

var (
	// defaultFields is the set of fields that are ignored by default for any
	// files by the given names.
	defaultFields = map[string][]string{
		"apply.json": {
			"0",
			"*.@timestamp",
		},
		"plan.json": {
			"terraform_version",
			"prior_state.terraform_version",
		},
		"state.json": {
			"terraform_version",
		},
	}
)

// TestOutput maps a Test case to a parsed set of JSON objects.
//
// The Files function will return these JSON objects, pre-stripped of any
// unwanted JSON fields.
type TestOutput struct {
	Test  Test
	files map[string]*files.File
}

// Files returns the JSON files that were returned by the test stripped of any
// unwanted fields.
func (output TestOutput) Files() (map[string]*files.File, error) {
	ret := map[string]*files.File{}
	for name, file := range output.files {
		contents, ok := file.Json()
		if !ok {
			// We only strip data out of JSON files, so skip this because it is
			// not a JSON file.
			ret[name] = file
			continue
		}

		var ignoreFields []string
		ignoreFields = append(ignoreFields, defaultFields[name]...)
		ignoreFields = append(ignoreFields, output.Test.Specification.IgnoreFields[name]...)

		stripped, err := strip.Strip(ignoreFields, contents)
		if err != nil {
			return nil, err
		}
		ret[name] = files.NewJsonFile(stripped)
	}
	return ret, nil
}

// ComputeDiff will report the difference between this TestOutput and the output
// already stored in the golden directory specified by the parameter.
func (output TestOutput) ComputeDiff(goldens string) (map[string]string, error) {
	newFiles, err := output.Files()
	if err != nil {
		return nil, err
	}

	ret := map[string]string{}
	for name, newFile := range newFiles {
		target := path.Join(goldens, output.Test.Name, name)

		goldenFile, err := os.ReadFile(target)
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, err
			}
		}

		if goldenFile == nil {
			// Then this means we don't have a golden file for this yet (as in
			// this is the first time we are using it). Let's just pretend it
			// was empty.
			ret[name] = NewFile
			continue
		}

		var diff string
		switch newFile.Ext() {
		case files.Json:
			// Then we can marshal the goldenFile into a JSON struct and get
			// more interesting output.
			newFileJson, _ := newFile.Json()
			var oldFileJson interface{}
			if err := json.Unmarshal(goldenFile, &oldFileJson); err != nil {
				return nil, err
			}
			diff = cmp.Diff(oldFileJson, newFileJson)
		case files.Raw:
			// Then we're just going to do a string comparison between the
			// goldenFile bytes and newFile.
			newFileString, _ := newFile.String()
			diff = cmp.Diff(string(goldenFile), newFileString)
		default:
			return nil, errors.New("found unrecognized file type: " + newFile.Ext())
		}

		if len(diff) == 0 {
			ret[name] = NoChange
		} else {
			ret[name] = diff
		}
	}
	return ret, nil
}

// UpdateGoldenFiles will write out the files for a given TestOutput into a
// target directory. This will overwrite any files already in the target
// directory.
func (output TestOutput) UpdateGoldenFiles(target string) error {
	tmp, err := os.MkdirTemp(target, output.Test.Name)
	if err != nil {
		return err
	}

	// We won't RemoveAll with tmp automatically, as there will be a point where
	// the original file has been deleted and tmp is all we have in which case
	// we don't want to delete tmp if anything goes wrong moving tmp into the
	// original location. tmp can be used by the user to recover manually.

	outputFiles, err := output.Files()
	if err != nil {
		return err
	}

	for name, file := range outputFiles {
		var data []byte
		switch file.Ext() {
		case files.Json:
			contents, _ := file.Json()
			var err error
			if data, err = json.MarshalIndent(contents, "", "  "); err != nil {
				os.RemoveAll(tmp)
				return err
			}
		case files.Raw:
			contents, _ := file.String()
			data = []byte(contents)
		}

		target := path.Join(tmp, name)
		if _, err := os.Stat(filepath.Dir(target)); os.IsNotExist(err) {
			// This means the parent directory for the target file doesn't exist
			// so let's make it.
			if err := os.MkdirAll(filepath.Dir(target), os.ModePerm); err != nil {
				os.RemoveAll(tmp)
				return err
			}
		}

		if err := os.WriteFile(target, data, os.ModePerm); err != nil {
			os.RemoveAll(tmp)
			return err
		}
	}

	// Now we've copied all the new golden files into our temporary directory,
	// we just need to move everything over to the original.
	if err := os.RemoveAll(path.Join(target, output.Test.Name)); err != nil {
		os.RemoveAll(tmp)
		return err
	}

	// From now, any failures are bad. We have removed the old directory so we
	// have lost the previous state of the golden files. If anything goes wrong
	// at this point we won't delete the tmp directory so that the user can
	// recover the failed test case manually by moving the tmp directory over
	// themselves.

	if err := os.Mkdir(path.Join(target, output.Test.Name), os.ModePerm); err != nil {
		return err
	}

	if err = filepath.WalkDir(tmp, files.CopyDir(tmp, path.Join(target, output.Test.Name), nil)); err != nil {
		return err
	}

	if err := os.RemoveAll(tmp); err != nil {
		return err
	}

	return nil
}
