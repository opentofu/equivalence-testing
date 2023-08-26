// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tests

import "github.com/opentffoundation/equivalence-testing/internal/binary"

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
	IncludeFiles []string                     `json:"include_files"`
	IgnoreFields map[string][]string          `json:"ignore_fields"`
	Rewrites     map[string]map[string]string `json:"rewrites"`

	// If Commands is empty, then we will execute a default set of commands:
	// [init, plan, apply, show, show plan]. Otherwise, these are the set of
	// commands that should be executed by the equivalence test framework for
	// this test case.
	Commands []binary.Command `json:"commands"`
}

func (s *TestSpecification) AddRewrites(rewrites map[string]map[string]string) {
	if s.Rewrites == nil {
		s.Rewrites = make(map[string]map[string]string)
	}

	for path, pathRewrites := range rewrites {
		current, exists := s.Rewrites[path]
		if !exists {
			s.Rewrites[path] = pathRewrites
		} else {
			for from, to := range pathRewrites {
				// Only add the rewrite if it doesn't already exist.
				if _, exists := s.Rewrites[path][from]; !exists {
					current[from] = to
				}
			}
		}
	}
}
