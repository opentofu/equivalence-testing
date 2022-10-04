package files

import (
	"encoding/json"
	"path/filepath"
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
		return f.contents.(string), true
	}
	return "", false
}
