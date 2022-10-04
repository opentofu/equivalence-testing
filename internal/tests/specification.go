package tests

import "github.com/hashicorp/terraform-equivalence-testing/internal/terraform"

// TestSpecification is a struct that provides the specification for a given
// test case.
//
// Each test has a set of additional files that should be included in the
// golden file update and diff functions, these are specified in the
// IncludeFiles field.
//
// Each test also has a set of JSON fields for each file that should be ignored
// when updating or diffing, these are specified in the IgnoreFields field.
type TestSpecification struct {
	IncludeFiles []string            `json:"include_files"`
	IgnoreFields map[string][]string `json:"ignore_fields"`

	// If Commands is empty, then we will execute a default set of commands:
	// [init, plan, apply, show, show plan]. Otherwise, these are the set of
	// commands that should be executed by the equivalence test framework for
	// this test case.
	Commands []terraform.Command `json:"commands"`
}
