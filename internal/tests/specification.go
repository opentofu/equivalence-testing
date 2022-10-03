package tests

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
}
