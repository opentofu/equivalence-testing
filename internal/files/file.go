// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package files

import (
	"encoding/json"
	"path/filepath"
	"strings"
)

const (
	Json = "json"
	Raw  = "raw"
)

func NewFile(file string, contents []byte) (*File, error) {
	if filepath.Ext(file) == ".json" {
		var data interface{}
		if err := json.Unmarshal(contents, &data); err != nil {
			return nil, err
		}
		return NewJsonFile(data), nil
	}
	return NewRawFile(string(contents)), nil
}

func NewRawFile(contents string) *File {
	return &File{
		contents: contents,
		ext:      Raw,
	}
}

func NewJsonFile(contents interface{}) *File {
	return &File{
		contents: contents,
		ext:      Json,
	}
}

type File struct {
	contents interface{}
	ext      string
	rewrites map[string]string
}

func (f File) Ext() string {
	return f.ext
}

func (f File) Json() (interface{}, bool) {
	if f.ext == Json {
		return f.contents, true
	}
	return nil, false
}

func (f File) String() (string, bool) {
	if f.ext == Raw {
		out := f.contents.(string)

		for from, to := range f.rewrites {
			out = strings.ReplaceAll(out, from, to)
		}

		return out, true
	}
	return "", false
}

func (f File) WithRewrites(rewrites map[string]string) *File {
	return &File{
		contents: f.contents,
		ext:      f.ext,
		rewrites: rewrites,
	}
}
